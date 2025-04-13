package internal

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestCheckErr(t *testing.T) {
	var buf bytes.Buffer
	exitCalled := false

	// Set up test logger
	Logger = zerolog.New(&buf).With().Timestamp().Logger()

	// Disable fatal to avoid os.Exit
	UseFatal = false

	// Override ExitFunc just in case
	ExitFunc = func(code int) {
		exitCalled = true
	}

	// Reset after test
	defer func() {
		Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		ExitFunc = os.Exit
		UseFatal = true
	}()

	CheckErr(errors.New("something went wrong"), "custom message")

	// Check the log output
	logged := buf.String()
	assert.Contains(t, logged, "custom message")
	assert.Contains(t, logged, "something went wrong")
	assert.False(t, exitCalled)
}
