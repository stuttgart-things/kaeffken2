package main

import (
	"context"
	"dagger/ci/internal/dagger"
)

func (m *Ci) Build(
	ctx context.Context,
	src *dagger.Directory) (
	*dagger.File,
	error) {

	// INITIALIZE THE GO MODULE
	goModule := dag.Go()

	// CALL THE BUILD FUNCTION WITH THE STRUCT
	buildOutput := goModule.Binary(
		src, // Source directory
		dagger.GoBinaryOpts{
			GoVersion:  goVersion,
			Os:         osVersion,
			Arch:       arch,
			GoMainFile: goMainFile,
			BinName:    binName,
			Ldflags:    ldflags,
		})

	binaryFile := buildOutput.File(binName)

	return binaryFile, nil
}
