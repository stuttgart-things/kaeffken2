package main

import (
	"context"

	"dagger/ci/internal/dagger"
)

var (
	goVersion  = "1.24.2"
	osVersion  = "linux"
	arch       = "amd64"
	goMainFile = "main.go"
	binName    = "kaeffken2"
	ldflags    = ""
	srcDir     = "."
)

type Ci struct {
	// +optional
	// +default="cgr.dev/chainguard/wolfi-base:latest"
	BaseImage string
}

func (m *Ci) container(
	ctx context.Context,
	src *dagger.Directory,
) *dagger.Container {

	if m.BaseImage == "" {
		m.BaseImage = "cgr.dev/chainguard/wolfi-base:latest"
	}

	ctr := dag.Container().From(m.BaseImage)

	// BUILD THE GO MODULE
	binaryFile, err := m.Build(ctx, src)
	if err != nil {
		panic(err)
	}

	return ctr.
		WithFile("/usr/bin/"+binName, binaryFile).
		WithExec([]string{"chmod", "+x", "/usr/bin/" + binName}).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")
}
