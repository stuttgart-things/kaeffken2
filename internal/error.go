package internal

import "log"

var FatalFunc = log.Fatalf // default

func CheckErr(err error, msg string) {
	if err != nil {
		FatalFunc("%s: %v", msg, err)
	}
}
