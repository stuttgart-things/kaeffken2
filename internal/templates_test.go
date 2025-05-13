package internal

import (
	"reflect"
	"testing"
)

func TestGetTemplatesPaths_WithTemplating(t *testing.T) {
	allAnswers := map[string]interface{}{
		"templates": []interface{}{
			"tests/{{ .cloud }}ansible-{{ .lab }}.k:/test/{{ .cloud }}ansible-{{ .lab }}.yaml",
		},
	}

	values := map[string]interface{}{
		"cloud": "vsphere",
		"lab":   "labul",
	}

	expected := []Template{
		{
			Source:      "tests/vsphereansible-labul.k",
			Destination: "/test/vsphereansible-labul.yaml",
		},
	}

	templates, err := GetTemplatesPaths(allAnswers, "templates", values)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(templates, expected) {
		t.Errorf("expected %v, got %v", expected, templates)
	}
}
