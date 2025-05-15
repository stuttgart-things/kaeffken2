/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceTripleQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple value",
			input:    `'''4'''`,
			expected: `'4'`,
		},
		{
			name:     "string value",
			input:    `'''helloopo'''`,
			expected: `'helloopo'`,
		},
		{
			name:     "with surrounding text",
			input:    `count: '''4'''`,
			expected: `count: '4'`,
		},
		{
			name:     "multiple replacements",
			input:    `count: '''4''', name: '''helloopo'''`,
			expected: `count: '4', name: 'helloopo'`,
		},
		{
			name:     "no replacement needed",
			input:    `count: '4'`,
			expected: `count: '4'`,
		},
		{
			name:     "empty value",
			input:    `count: ''''''`,
			expected: `count: ''`,
		},
		{
			name:     "mixed quotes",
			input:    `count: '''4''', name: "hello"`,
			expected: `count: '4', name: "hello"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceTripleQuotes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFixQuotesInMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name: "simple values",
			input: map[string]string{
				"count": "'''4'''",
				"name":  "'''helloopo'''",
			},
			expected: map[string]string{
				"count": "'4'",
				"name":  "'helloopo'",
			},
		},
		{
			name: "mixed quotes",
			input: map[string]string{
				"cpu":   "'''8'''",
				"ram":   "'''4'''",
				"disk":  "'1'",    // already single quoted
				"label": `"fast"`, // double quoted
			},
			expected: map[string]string{
				"cpu":   "'8'",
				"ram":   "'4'",
				"disk":  "'1'",
				"label": `"fast"`,
			},
		},
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixQuotesInMap(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("fixQuotesInMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("malformed triple quotes", func(t *testing.T) {
		input := `'''unmatched`
		result := replaceTripleQuotes(input)
		assert.Equal(t, input, result) // Should remain unchanged
	})

	t.Run("nested quotes", func(t *testing.T) {
		input := `'''"value"'''`
		result := replaceTripleQuotes(input)
		assert.Equal(t, `'"value"'`, result)
	})

	t.Run("real world example", func(t *testing.T) {
		input := `config: {count: '''4''', name: '''hello''', nested: {value: '''test'''}}`
		expected := `config: {count: '4', name: 'hello', nested: {value: 'test'}}`
		result := replaceTripleQuotes(input)
		assert.Equal(t, expected, result)
	})
}
