/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package modules

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func GetSpec(specName string) (map[string]interface{}, error) {
	// Define the YAML data as a string
	yamlData := `
specs:
  small:
    cpu: 1
    memory: 2Gi
    disk: 20Gi
  medium:
    cpu: 2
    memory: 4Gi
    disk: 40Gi
  large:
    cpu: 4
    memory: 8Gi
    disk: 80Gi
  xlarge:
    cpu: 8
    memory: 16Gi
    disk: 160Gi
`

	// A generic map to hold the unmarshalled YAML data
	var data map[string]map[string]map[string]interface{}

	// Unmarshal the YAML into the generic structure
	err := yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	// Extract the spec section (which is under the `specs` key)
	specs, ok := data["specs"]
	if !ok {
		return nil, fmt.Errorf("specs not found in YAML")
	}

	// Retrieve the requested spec (small, medium, large, xlarge)
	spec, ok := specs[specName]
	if !ok {
		return nil, fmt.Errorf("spec %s not found", specName)
	}

	return spec, nil
}
