package main

import "context"

// Test runs the Go tests
func (m *GoDagger) Test(ctx context.Context,
	// source is the directory containing the Go source code
	// +required
	source *Directory,
	// goVersion is the version of Go to use for building the binary
	// +optional
	// +default="1.22.0"
	goVersion string,
	// ignoreCache is a flag to ignore the cache
	// +optional
	// +default=false
	ignoreCache bool,
	// race is a flag to include Go's built-in data race detector
	// +optional
	// +default=false
	race bool,
	// verbose is a flag to enable verbose output
	// +optional
	// +default=false
	verbose bool,
) (*Container, error) {
	cli := dag.Pipeline("go-tests")

	ignoreCacheFlag := ""
	if ignoreCache {
		ignoreCacheFlag = "--count=1"
	}

	raceFlag := ""
	if race {
		raceFlag = "-race"
	}

	verboseFlag := ""
	if verbose {
		verboseFlag = "-v"
	}

	return cli.Container().
		From("golang:"+goVersion).
		WithDirectory("/src", source).
		WithMountedCache("/go/pkg/mod", m.goModCacheVolume()).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", verboseFlag, raceFlag, ignoreCacheFlag, "./..."}).
		Sync(ctx)
}
