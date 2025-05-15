/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"bytes"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestCheckErr(t *testing.T) {
	var buf bytes.Buffer
	exitCalled := false

	// INIT LOGGER
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	// Disable fatal to avoid os.Exit
	UseFatal = false

	// Override ExitFunc just in case
	ExitFunc = func(code int) {
		exitCalled = true
	}

	// Reset after test
	defer func() {
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
