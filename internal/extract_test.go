package internal

import (
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
