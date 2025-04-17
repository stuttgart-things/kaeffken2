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
