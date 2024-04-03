package configloader

import "github.com/snippetaccumulator/configloader/fieldsetter"

type MockLoader struct {
	MockData  map[string]any
	overrides map[string]any
}

// NewMockLoader creates an instance of MockLoader with predefined mock data. This mock loader is designed
// for testing purposes, allowing developers to simulate the loading of configuration data without accessing
// actual files. The mock data should be structured as a map representing the configuration fields and their
// desired values, facilitating easy and controlled testing scenarios.
func NewMockLoader(mockData map[string]any) *MockLoader {
	return &MockLoader{
		MockData:  mockData,
		overrides: make(map[string]any),
	}
}

// Load simulates the configuration loading process using predefined mock data. This method attempts to set
// the configuration fields of the provided config object based on the mock data and any overrides that have
// been specified. It is designed for testing, allowing developers to verify the behavior of their applications
// with various configurations without needing to interact with actual configuration files.
// Returns the first error encountered during field setting if any, otherwise nil.
func (m *MockLoader) Load(config any) error {
	errs := fieldsetter.SetFields(config, m.MockData, false)
	if len(errs) > 0 {
		return errs[0]
	}

	errs = fieldsetter.SetFields(config, m.overrides, true)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Override allows for the specification of override values in the mock configuration loading process. This method
// adds or updates entries in the mock loader's internal overrides map, enabling the testing of configuration behavior
// when specific fields are modified after the initial mock data load. It's particularly useful for validating
// how applications handle dynamic configuration changes.
func (m *MockLoader) Override(path string, value any) error {
	m.overrides[path] = value
	return nil
}
