package internal

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckErr(t *testing.T) {
	called := false
	var gotMsg string

	// Mock FatalFunc to capture the final message
	FatalFunc = func(format string, v ...any) {
		called = true
		gotMsg = fmt.Sprintf(format, v...)
	}

	defer func() { FatalFunc = nil }() // Reset after test

	CheckErr(errors.New("something went wrong"), "custom message")

	assert.True(t, called)
	assert.Contains(t, gotMsg, "custom message")
	assert.Contains(t, gotMsg, "something went wrong")
}
