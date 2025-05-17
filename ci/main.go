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

func (m *Ci) Debug(
	ctx context.Context,
	// +ignore=["*", "!analytics"]
	source *dagger.Directory,
) *dagger.Directory {
	return source
}

func (m *Ci) TestNonInteractive(
	ctx context.Context,
	// "**", "!**"
	src *dagger.Directory) {

	// LIST ALL ENTRIES
	entries, err := src.Entries(ctx)
	if err != nil {
		panic(err)
	}

	// PRINT ALL ENTRIES
	for _, entry := range entries {
		println(entry)
	}

	// BUILD THE GO MODULE
	buildOutput, err := m.Build(ctx, src)

	if err != nil {
		panic(err)
	}

	// Extract the binary file from the build output directory
	binaryFile := buildOutput.File(binName)

	testContainer := m.container().
		WithFile("/usr/bin/"+binName, binaryFile).
		WithExec([]string{"chmod", "+x", "/usr/bin/" + binName}).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")

	// PRINT SRC DIRECTORY
	_, err = testContainer.
		WithExec([]string{"tree"}).Stdout(ctx)

	// PRINT VERSION
	_, err = testContainer.
		WithExec([]string{binName, "version"}).Stdout(ctx)

	_, err = testContainer.
		WithExec([]string{binName,
			"render",
			"--config", "tests/vmRequestConfig.yaml",
			"--request", "tests/vmRequest.yaml",
			"--survey=false",
		}).Stdout(ctx)

	if err != nil {
		panic(err)
	}

}

func (m *Ci) Build(
	ctx context.Context,
	src *dagger.Directory) (*dagger.Directory, error) {

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

	return buildOutput, nil

}

func (m *Ci) container() *dagger.Container {
	if m.BaseImage == "" {
		m.BaseImage = "cgr.dev/chainguard/wolfi-base:latest"
	}

	ctr := dag.Container().From(m.BaseImage)

	return ctr
}
