/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"reflect"
	"testing"
)

func TestMergeMaps(t *testing.T) {
	a := map[string]interface{}{
		"one": 1,
		"two": "second",
	}

	b := map[string]interface{}{
		"two":   "overwritten",
		"three": 3,
	}

	expected := map[string]interface{}{
		"one":   1,
		"two":   "overwritten",
		"three": 3,
	}

	result := MergeMaps(a, b)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCleanMap(t *testing.T) {
	input := map[string]interface{}{
		"name":  "example",
		"empty": "",
		"nil":   nil,
		"slice": []interface{}{},
		"nested": map[string]interface{}{
			"foo": "",
			"bar": "baz",
		},
	}

	expected := map[string]interface{}{
		"name": "example",
		"nested": map[string]interface{}{
			"bar": "baz",
		},
	}

	result := CleanMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("CleanMap() = %#v, want %#v", result, expected)
	}
}
