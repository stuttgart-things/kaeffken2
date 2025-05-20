package main

import (
	"context"
	"dagger/ci/internal/dagger"
)

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
