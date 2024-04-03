package configloader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/snippetaccumulator/configloader/fieldsetter"
)

// ConfigLoader is a struct designed to load and apply configurations from various sources. It supports loading
// main configuration and override configuration from specified paths and filenames, applying deserializers for
// each configuration format, and dynamically overriding specific configuration fields via a map of paths to values.
// Fields include Name, Path, OverrideName, OverridePath for file locations, Deserializer, and OverrideDeserializer
// for handling specific data formats, and Overrides for field-specific overrides.
type ConfigLoader struct {
	Name                 string
	Path                 string
	OverrideName         string
	OverridePath         string
	Deserializer         DeserializerFunc
	OverrideDeserializer DeserializerFunc
	Overrides            map[string]any
}

// NewConfigLoader creates and returns a new instance of ConfigLoader with the specified name. It initializes
// the ConfigLoader's fields with default values: current directory for Path, empty for OverrideName and
// OverridePath, nil for Deserializer and OverrideDeserializer, and an empty map for Overrides. Additional
// configurations can be applied using Option functions passed as arguments to this function, allowing for
// customization of the loader's behavior and settings.
func NewConfigLoader(name string, options ...Option) *ConfigLoader {
	loader := &ConfigLoader{
		Name:         name,
		Path:         ".",
		OverrideName: "",
		OverridePath: "",
		Deserializer: nil,
		Overrides:    make(map[string]any),
	}
	for _, option := range options {
		option(loader)
	}
	return loader
}

// Load reads the main configuration file based on the ConfigLoader's Path and Name, deserializes it into
// the provided config object using the set Deserializer, and applies any Overrides. If OverridePath and
// OverrideName are set, it also loads and applies an override configuration file using either the OverrideDeserializer
// or the main Deserializer if no OverrideDeserializer is set. Errors during file reading, deserialization, or
// field setting are returned. This method facilitates the flexible loading and merging of configurations with
// optional overrides to tailor application settings dynamically.
func (c *ConfigLoader) Load(config any) error {
	if c.Deserializer == nil {
		return fmt.Errorf("no deserializer set for main configuration")
	}

	configData, err := os.ReadFile(filepath.Join(c.Path, c.Name))
	if err != nil {
		return err
	}

	err = c.Deserializer.Deserialize(configData, config)
	if err != nil {
		return err
	}

	if c.OverrideName != "" && c.OverridePath != "" {
		if c.OverrideDeserializer == nil {
			c.OverrideDeserializer = c.Deserializer
		}
		overrideData, err := os.ReadFile(filepath.Join(c.OverridePath, c.OverrideName))
		if err != nil {
			return err
		}

		err = c.OverrideDeserializer.Deserialize(overrideData, config)
		if err != nil {
			return err
		}
	}

	errs := fieldsetter.SetFields(config, c.Overrides, true)
	if len(errs) > 0 {
		return fmt.Errorf("error setting fields: %+v", errs)
	}

	return nil
}

// Override adds or updates a specific configuration override by path. The path should specify the target field
// within the configuration object, and the value is what will be set for this field when applying overrides.
// This method allows for dynamic adjustments to the configuration, even after the initial loading process.
func (c *ConfigLoader) Override(path string, value any) error {
	c.Overrides[path] = value
	return nil
}
