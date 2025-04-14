package internal

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stuttgart-things/survey"
)

func TestParseKCLQuestions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []*survey.Question
		wantErr  bool
	}{
		{
			name:  "basic ask question",
			input: `_name = option("name") or "" # ask;`,
			expected: []*survey.Question{
				{
					Name:      "name",
					Default:   "",
					Kind:      "ask",
					Type:      "string",
					Prompt:    "NAME?",
					Options:   nil,
					MinLength: 0,
					MaxLength: 0,
				},
			},
			wantErr: false,
		},
		{
			name:  "ask question with default value",
			input: `_cpu = option("cpu") or "8" # ask;`,
			expected: []*survey.Question{
				{
					Name:      "cpu",
					Default:   "8",
					Kind:      "ask",
					Type:      "string",
					Prompt:    "CPU?",
					Options:   nil,
					MinLength: 0,
					MaxLength: 0,
				},
			},
			wantErr: false,
		},
		{
			name:  "select question with options",
			input: `_cpu = option("cpu") or "8" #select;4,8,12,16`,
			expected: []*survey.Question{
				{
					Name:      "cpu",
					Default:   "8",
					Kind:      "select",
					Prompt:    "CPU?",
					Options:   []string{"4", "8", "12", "16"},
					MinLength: 0,
					MaxLength: 0,
				},
			},
			wantErr: false,
		},
		{
			name:  "question with length constraints",
			input: `_name = option("name") or "" # ask;4,8,12,16-min2+max30`,
			expected: []*survey.Question{
				{
					Name:      "name",
					Default:   "",
					Kind:      "ask",
					Type:      "string",
					Prompt:    "NAME?",
					Options:   nil, // Note: Options should only be set for "select" type
					MinLength: 2,
					MaxLength: 30,
				},
			},
			wantErr: false,
		},
		{
			name: "multiple questions",
			input: `_name = option("name") or "" # ask;4,8,12,16-min2+max30
_cpu = option("cpu") or "8" #select;4,8,12,16`,
			expected: []*survey.Question{
				{
					Name:      "name",
					Default:   "",
					Kind:      "ask",
					Type:      "string",
					Prompt:    "NAME?",
					Options:   nil,
					MinLength: 2,
					MaxLength: 30,
				},
				{
					Name:      "cpu",
					Default:   "8",
					Kind:      "select",
					Prompt:    "CPU?",
					Options:   []string{"4", "8", "12", "16"},
					MinLength: 0,
					MaxLength: 0,
				},
			},
			wantErr: false,
		},
		{
			name:     "no match",
			input:    `_name = "static value" # not a question`,
			expected: nil, // Changed to expect nil
			wantErr:  false,
		},
		{
			name:     "malformed line",
			input:    `_name = option("name" or "" # ask;`,
			expected: nil, // Changed to expect nil
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseKCLQuestions(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseKCLQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestExtractQuestionsFromKCLFile(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		_, err := ExtractQuestionsFromKCLFile("nonexistent.k")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read KCL file")
	})
}

func TestParseTemplateValues(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect map[string]interface{}
	}{
		{
			name:  "Basic string",
			input: []string{"name=web"},
			expect: map[string]interface{}{
				"name": "web",
			},
		},
		{
			name:  "Integer and boolean values",
			input: []string{"replicas=3", "debug=true"},
			expect: map[string]interface{}{
				"replicas": 3,
				"debug":    true,
			},
		},
		{
			name:  "Invalid format is ignored",
			input: []string{"invalid", "key=value"},
			expect: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:  "Multiple mixed types",
			input: []string{"stringVal=hello", "intVal=42", "boolVal=false"},
			expect: map[string]interface{}{
				"stringVal": "hello",
				"intVal":    42,
				"boolVal":   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseTemplateValues(tt.input)
			if !reflect.DeepEqual(actual, tt.expect) {
				t.Errorf("Expected %v, got %v", tt.expect, actual)
			}
		})
	}
}

func TestReadDictEntry(t *testing.T) {
	yamlContent := `
dicts:
  kinds:
    labul_proxmoxvm:
      env: labul
      cloud: proxmox
      kind: proxmoxvmansible
      template: proxmoxvmansible-labul.k
    labda_vspherevm:
      env: labda
      cloud: vsphere
      kind: vspherevmansible
      template: vspherevmansible-labda.k
`

	tmpFile, err := os.CreateTemp("", "dicts_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(yamlContent); err != nil {
		t.Fatalf("Failed to write YAML to temp file: %v", err)
	}
	tmpFile.Close()

	expected := map[string]interface{}{
		"env":      "labda",
		"cloud":    "vsphere",
		"kind":     "vspherevmansible",
		"template": "vspherevmansible-labda.k",
	}

	result, err := ReadDictEntry(tmpFile.Name(), "kinds", "labda_vspherevm")
	if err != nil {
		t.Fatalf("ReadDictEntry failed: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %+v, Got: %+v", expected, result)
	}
}

func TestReadDicts(t *testing.T) {
	yamlContent := `
dicts:
  kinds:
    labul_proxmoxvm:
      env: labul
      cloud: proxmox
      kind: proxmoxvmansible
      template: proxmoxvmansible-labul.k
  sizes:
    small:
      cpu: 1
      memory: 2Gi
      disk: 20Gi
`

	// Create a temporary YAML file
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	// Read the dicts from the file
	dicts, err := ReadDicts(yamlPath, "dicts")
	assert.NoError(t, err)
	assert.NotNil(t, dicts)

	// Assert top-level keys
	_, ok := dicts["kinds"]
	assert.True(t, ok, "'kinds' key should be present")

	_, ok = dicts["sizes"]
	assert.True(t, ok, "'sizes' key should be present")

	// Assert nested structure
	kinds := dicts["kinds"].(map[string]interface{})
	_, ok = kinds["labul_proxmoxvm"]
	assert.True(t, ok, "'labul_proxmoxvm' should be under 'kinds'")
}
