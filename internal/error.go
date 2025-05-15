/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ExitFunc = os.Exit
	UseFatal = true // Allows tests to avoid fatal exit
)

func CheckErr(err error, msg string) {

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	if err != nil {
		event := log.Error()
		if UseFatal {
			event = log.Fatal()
		}

		event.
			Err(err).
			Str("context", msg).
			Msg("fatal error")

		if UseFatal {
			ExitFunc(1)
		}
	}
}
