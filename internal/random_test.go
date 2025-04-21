package internal

import (
	"testing"
)

func TestGetRandomStringFromMap(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		key       string
		wantError bool
	}{
		{
			name: "Valid key with string slice",
			input: map[string]interface{}{
				"colors": []interface{}{"red", "green", "blue"},
			},
			key:       "colors",
			wantError: false,
		},
		{
			name: "Missing key",
			input: map[string]interface{}{
				"animals": []interface{}{"cat", "dog"},
			},
			key:       "fruits",
			wantError: true,
		},
		{
			name: "Wrong type at key",
			input: map[string]interface{}{
				"numbers": 123,
			},
			key:       "numbers",
			wantError: true,
		},
		{
			name: "Empty slice",
			input: map[string]interface{}{
				"empty": []interface{}{},
			},
			key:       "empty",
			wantError: true,
		},
		{
			name: "Slice with no strings",
			input: map[string]interface{}{
				"mixed": []interface{}{1, 2, 3},
			},
			key:       "mixed",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRandomStringFromMap(tt.input, tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("getRandomStringFromMap() error = %v, wantError = %v", err, tt.wantError)
			}
			if err == nil && got == "" {
				t.Errorf("Expected non-empty string, got empty")
			}
		})
	}
}
