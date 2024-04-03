package configloader

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	env "github.com/Netflix/go-env"
	"gopkg.in/yaml.v3"
)

// JSONDeserializer is a struct with no fields, implementing the DeserializerFunc interface for JSON data.
// It provides a method to deserialize JSON encoded data into a Go value.
type JSONDeserializer struct{}

func (jd *JSONDeserializer) Deserialize(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// YAMLDeserializer is a struct with no fields, implementing the DeserializerFunc interface for YAML data.
// It offers a method to deserialize YAML encoded data into a Go value.
type YAMLDeserializer struct{}

func (yd *YAMLDeserializer) Deserialize(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

// TOMLDeserializer implements the DeserializerFunc interface for TOML data.
// It provides a method to deserialize TOML encoded data into a Go value.
type TOMLDeserializer struct{}

func (td *TOMLDeserializer) Deserialize(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}

// EnvDeserializer implements the DeserializerFunc interface for environment variable data like .env files.
type EnvDeserializer struct{}

func (ed *EnvDeserializer) Deserialize(_ []byte, v any) error {
	_, err := env.UnmarshalFromEnviron(v)
	return err
}
