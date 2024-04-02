package main

import "context"

// Lint runs the Go linter
func (m *GoDagger) Lint(ctx context.Context,
	// source is the directory containing the Go source code
	// +required
	source *Directory,
	// golangCILintImage is the container image to use for the linter
	// +optional
	// +default="golangci/golangci-lint:v1.57.2"
	golangCILintImage string,
	// timeout is the maximum time to run the linter
	// +optional
	// +default="5m"
	timeout string,
	// verbose is a flag to enable verbose output
	// +optional
	// +default=false
	verbose bool,
) (*Container, error) {
	cli := dag.Pipeline("go-lint")

	args := []string{"golangci-lint", "run", "--timeout", timeout}

	if verbose {
		args = append(args, "-v")
	}

	return cli.Container().From(golangCILintImage).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Sync(ctx)
}
