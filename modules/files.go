/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package modules

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/stuttgart-things/kaeffken2/internal"
)

var (
	err error
)

type InputFile struct {
	Name       string
	Path       string
	FileExists bool
}

func CheckInputFiles(inputFiles []InputFile) (configSpec, configValues, requestSpec map[string]interface{}) {

	// CHECK IF GIVEN FILES EXISTS
	for i := range inputFiles {
		exists, err := internal.FileExists(inputFiles[i].Path)
		if err != nil {
			log.Error().Err(err).Str("path", inputFiles[i].Path).Msg("Error checking file")
		}
		inputFiles[i].FileExists = exists
	}

	// DECLARE VARIABLES
	configValues = make(map[string]interface{})
	configSpec = make(map[string]interface{})
	requestSpec = make(map[string]interface{})

	for _, f := range inputFiles {
		switch f.Name + fmt.Sprintf(":%t", f.FileExists) {
		case "template:true":
			log.Info().Str("path", f.Path).Msg("Template exists ✅")

		case "template:false":
			log.Warn().Str("path", f.Path).Msg("Template missing ❌")

		case "request:true":
			log.Info().Str("path", f.Path).Msg("Request exists ✅")

			// READ REQUEST SPEC
			requestSpec, err = internal.ReadSpecSection(f.Path)
			internal.CheckErr(err, "ERROR READING REQUEST SPEC")

		case "request:false":
			log.Warn().Str("path", f.Path).Msg("Request missing ❌")

		case "config:true":
			log.Info().Str("path", f.Path).Msg("Config exists ✅")

			// READ CONFIG SPEC
			configSpec, err = internal.ReadSpecSection(f.Path)
			internal.CheckErr(err, "ERROR READING CONFIG SPEC")

			// READ CONFIG VALUES
			configValues, err = internal.ReadDicts(f.Path, "dicts")
			internal.CheckErr(err, "ERROR READING CONFIG DICTS")

		case "config:false":
			log.Warn().Str("path", f.Path).Msg("Config missing ❌")

		default:
			log.Warn().Str("name", f.Name).Str("path", f.Path).Msg("Unknown input file type or state ❌")
		}
	}

	return configSpec, configValues, requestSpec
}
