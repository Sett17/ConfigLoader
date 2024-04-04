package configloader_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/snippetaccumulator/configloader"
)

type Config struct {
	Field1 string `yaml:"field1"`
	Field2 int    `yaml:"field2"`
	Nested struct {
		Field3 bool `yaml:"field3"`
	}
}

// Helper function to create a temporary YAML file
func createTempYAMLFile(content []byte) (string, error) {
	tmpfile, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		return "", err
	}
	if _, err := tmpfile.Write(content); err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		os.Remove(tmpfile.Name())
		return "", err
	}
	return tmpfile.Name(), nil
}

func TestLoadWithYAMLDeserializer(t *testing.T) {
	configData := []byte("field1: value1\nfield2: 2\nnested:\n  field3: true")
	filename, err := createTempYAMLFile(configData)
	if err != nil {
		t.Fatalf("Unable to create temp YAML file: %s", err)
	}
	defer os.Remove(filename)

	// Print the absolute path of the temp file for debugging
	absPath, _ := filepath.Abs(filename)
	fmt.Println("Temp file path:", absPath)

	var config Config
	loader := configloader.NewConfigLoader(filepath.Base(filename), configloader.WithPath(filepath.Dir(absPath)), configloader.WithDeserializer(new(configloader.YAMLDeserializer)))

	if err := loader.Load(&config); err != nil {
		t.Fatalf("Failed to load configuration: %s", err)
	}
}

func TestDeserializerNotSet(t *testing.T) {
	var config Config
	loader := configloader.NewConfigLoader("nonexistent.yaml", configloader.WithPath("."))

	err := loader.Load(&config)
	if err == nil {
		t.Error("Expected an error when deserializer is not set, got nil")
	}
}

func TestWithOverrideFile(t *testing.T) {
	mainConfigData := []byte("field1: main1\nfield2: 1\nnested:\n  field3: false")
	mainFilename, err := createTempYAMLFile(mainConfigData)
	if err != nil {
		t.Fatalf("Unable to create main temp YAML file: %s", err)
	}
	defer os.Remove(mainFilename)

	overrideConfigData := []byte("field1: override1\nnested:\n  field3: true")
	overrideFilename, err := createTempYAMLFile(overrideConfigData)
	if err != nil {
		t.Fatalf("Unable to create override temp YAML file: %s", err)
	}
	defer os.Remove(overrideFilename)

	var config Config
	loader := configloader.NewConfigLoader(filepath.Base(mainFilename),
		configloader.WithPath(filepath.Dir(mainFilename)),
		configloader.WithDeserializer(new(configloader.YAMLDeserializer)),
		configloader.WithOverrideFile(filepath.Dir(overrideFilename), filepath.Base(overrideFilename)),
	)

	if err := loader.Load(&config); err != nil {
		t.Fatalf("Failed to load configuration with override: %s", err)
	}

	if config.Field1 != "override1" {
		t.Errorf("Expected field1 to be 'override1', got '%s'", config.Field1)
	}
}

func TestOverrideValue(t *testing.T) {
	configData := []byte("field1: value1\nfield2: 2")
	filename, err := createTempYAMLFile(configData)
	if err != nil {
		t.Fatalf("Unable to create temp YAML file: %s", err)
	}
	defer os.Remove(filename)

	var config Config
	loader := configloader.NewConfigLoader(filepath.Base(filename),
		configloader.WithPath(filepath.Dir(filename)),
		configloader.WithDeserializer(new(configloader.YAMLDeserializer)),
	)

	// Override a value before loading
	err = loader.Override("Field1", "overridden")
	if err != nil {
		t.Fatalf("Failed to set override: %s", err)
	}
	err = loader.Override("Nested.Field3", true)
	if err != nil {
		t.Fatalf("Failed to set override: %s", err)
	}

	if err := loader.Load(&config); err != nil {
		t.Fatalf("Failed to load configuration: %s", err)
	}

	if config.Field1 != "overridden" {
		t.Errorf("Expected 'field1' to be overridden with 'overridden', got '%s'", config.Field1)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	var config Config
	loader := configloader.NewConfigLoader("does_not_exist.yaml",
		configloader.WithPath("."),
		configloader.WithDeserializer(new(configloader.YAMLDeserializer)),
	)

	err := loader.Load(&config)
	if err == nil {
		t.Error("Expected an error for non-existent file, got nil")
	}
}

type ConfigWithEmbeds struct {
	Config
	Field4 float64 `yaml:"field4"`
}

func TestLoadWithEmbeddedStructIndirectly(t *testing.T) {
	configData := []byte("config:\n  field1: value1\n  field2: 2\n  nested:\n    field3: true\nfield4: 3.14")
	filename, err := createTempYAMLFile(configData)
	if err != nil {
		t.Fatalf("Unable to create temp YAML file: %s", err)
	}

	var config ConfigWithEmbeds
	loader := configloader.NewConfigLoader(filepath.Base(filename),
		configloader.WithPath(filepath.Dir(filename)),
		configloader.WithDeserializer(new(configloader.YAMLDeserializer)),
	)

	if err := loader.Load(&config); err != nil {
		t.Fatalf("Failed to load configuration: %s", err)
	}

	if config.Field1 != "value1" {
		t.Errorf("Expected 'field1' to be 'value1', got '%s'", config.Field1)
	}
	if config.Field2 != 2 {
		t.Errorf("Expected 'field2' to be 2, got %d", config.Field2)
	}
	if config.Nested.Field3 != true {
		t.Errorf("Expected 'nested.field3' to be true, got %t", config.Nested.Field3)
	}
	if config.Field4 != 3.14 {
		t.Errorf("Expected 'field4' to be 3.14, got %f", config.Field4)
	}
}

func TestOverrideWithEmbeddedStruct(t *testing.T) {
	configData := []byte("config:\n  field1: value1\n  field2: 2\n  nested:\n    field3: true\nfield4: 3.14")
	filename, err := createTempYAMLFile(configData)
	if err != nil {
		t.Fatalf("Unable to create temp YAML file: %s", err)
	}

	var config ConfigWithEmbeds
	loader := configloader.NewConfigLoader(filepath.Base(filename),
		configloader.WithPath(filepath.Dir(filename)),
		configloader.WithDeserializer(new(configloader.YAMLDeserializer)),
	)

	loader.Override("Field1", "overridden")
	loader.Override("Config.Field2", 3)

	if err := loader.Load(&config); err != nil {
		t.Fatalf("Failed to load configuration: %s", err)
	}

	if config.Field1 != "overridden" {
		t.Errorf("Expected 'field1' to be 'overridden', got '%s'", config.Field1)
	}
	if config.Field2 != 3 {
		t.Errorf("Expected 'field2' to be 3, got %d", config.Field2)
	}
	if config.Nested.Field3 != true {
		t.Errorf("Expected 'nested.field3' to be true, got %t", config.Nested.Field3)
	}
	if config.Field4 != 3.14 {
		t.Errorf("Expected 'field4' to be 3.14, got %f", config.Field4)
	}
}
