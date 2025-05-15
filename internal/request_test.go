/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadSpecSection(t *testing.T) {
	// Sample YAML content
	yamlContent := `
apiVersion: resources.stuttgart-things.com/v1alpha1
kind: vmRequest
metadata:
  name: test-vm
spec:
  kind: labul_proxmoxvm
  size: large
  tags:
    - dev
    - test
`

	// Create a temporary YAML file
	tmpFile, err := os.CreateTemp("", "test-spec-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yamlContent))
	assert.NoError(t, err)
	tmpFile.Close()

	// Call the function
	spec, err := ReadSpecSection(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, spec)

	// Check some expected values
	assert.Equal(t, "labul_proxmoxvm", spec["kind"])
	assert.Equal(t, "large", spec["size"])
	assert.ElementsMatch(t, []interface{}{"dev", "test"}, spec["tags"].([]interface{}))
}
