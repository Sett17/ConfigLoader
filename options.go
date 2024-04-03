package configloader

// Option defines a function signature for optional configuration functions that customize the behavior of a ConfigLoader instance.
// These functions enable flexible and modular configuration of a ConfigLoader by setting various parameters such as file paths,
// deserializers, and override mechanisms. Each option function accepts a pointer to a ConfigLoader instance and modifies it
// according to the specific needs of the application setup.
type Option func(loader *ConfigLoader)

// WithPath configures the path where the main configuration file is located. This option function is intended
// to be used when creating a new ConfigLoader instance, allowing the specification of a custom directory
// path for the configuration file, overriding the default current directory.
func WithPath(path string) Option {
	return func(loader *ConfigLoader) {
		loader.Path = path
	}
}

// WithOverrideFile sets the path and name for an override configuration file. This option function enables
// the ConfigLoader to apply additional configuration settings from a specified file, allowing for flexible
// adjustments and customization beyond the main configuration. The override file is processed after the
// main configuration file, with its settings potentially overriding those set by the main configuration.
func WithOverrideFile(path, name string) Option {
	return func(loader *ConfigLoader) {
		loader.OverridePath = path
		loader.OverrideName = name
	}
}

// WithDeserializer is an option function for ConfigLoader that sets the specified deserializer function
// for interpreting the main configuration file. This allows for custom deserialization logic to be applied,
// enabling the support of various data formats beyond the default ones provided.
func WithDeserializer(deserializer DeserializerFunc) Option {
	return func(loader *ConfigLoader) {
		loader.Deserializer = deserializer
	}
}

// WithOverrideDeserializer specifies the deserializer function for the override configuration file.
// This option allows the ConfigLoader to interpret and apply settings from the override file using
// a custom deserialization logic, providing the flexibility to handle different data formats for
// override configurations.
func WithOverrideDeserializer(deserializer DeserializerFunc) Option {
	return func(loader *ConfigLoader) {
		loader.OverrideDeserializer = deserializer
	}
}
