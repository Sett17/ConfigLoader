package configloader

// Loader defines the interface for loading configuration data into a Go struct. It requires implementation
// of the Load method, which takes a configuration object and populates it with data, and the Override method,
// which allows for dynamic modifications of specific configuration fields after the initial load. This interface
// is designed to abstract the configuration loading mechanism, enabling the use of various sources and formats.
type Loader interface {
	Load(config any) error
	Override(path string, value any) error
}

// DeserializerFunc outlines the interface for deserializing data from a byte slice into a Go struct. It
// necessitates an implementation of the Deserialize method, which should interpret the provided byte slice
// according to the deserializer's format (e.g., JSON, YAML, TOML, etc.) and populate the passed-in struct
// accordingly. This interface supports flexible data interpretation, facilitating the integration of different
// configuration file formats.
type DeserializerFunc interface {
	Deserialize(data []byte, v any) error
}
