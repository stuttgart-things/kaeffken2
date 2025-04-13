package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ReadSpecSection loads the YAML file and returns the `spec` section as a generic map.
func ReadSpecSection(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	spec, ok := raw["spec"].(map[string]interface{})
	if !ok {
		return nil, nil // or return error if required
	}

	return spec, nil
}
