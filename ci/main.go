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

func (m *Ci) TestRenderCommand(
	ctx context.Context,
	// "**", "!**"
	src *dagger.Directory) {

	// LIST ALL FILE ENTRIES
	entries, err := src.Entries(ctx)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		println(entry)
	}

	// CREATE TESTE CONTAINER
	testContainer := m.container(ctx, src)

	// PRINT SRC DIRECTORY
	_, err = testContainer.
		WithExec([]string{"tree"}).Stdout(ctx)

	// PRINT VERSION
	_, err = testContainer.
		WithExec([]string{binName, "version"}).Stdout(ctx)

	// TEST RENDER COMMAND
	renderTest := testContainer.
		WithExec([]string{binName,
			"render",
			"--config", "tests/vmRequestConfig.yaml",
			"--request", "tests/vmRequest.yaml",
			"--survey=false",
		})

	_, err = renderTest.
		WithExec([]string{"tree", "/tmp"}).Stdout(ctx)

	if err != nil {
		panic(err)
	}

}

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
