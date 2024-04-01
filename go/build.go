package main

import (
	"context"
	"dagger/go-dagger/internal/dagger"
	"fmt"
	"path/filepath"

	platformFormat "github.com/containerd/containerd/platforms"
)

// Build builds the Go binary for the specified go version and platforms
func (m *GoDagger) Build(ctx context.Context,
	// source is the directory containing the Go source code
	// +required
	source *Directory,
	// goVersion is the version of Go to use for building the binary
	// +optional
	// +default="1.22.0"
	goVersion string,
	// platforms is the list of platforms to build the binary for
	// +optional
	// +default=["linux/amd64", "linux/arm64"]
	platforms []string) (*Directory, error) {
	cli := dag.Pipeline("go-build")

	var files []*File

	for _, platform := range platforms {
		os := platformFormat.MustParse(string(platform)).OS
		arch := platformFormat.MustParse(string(platform)).Architecture
		binaryName := fmt.Sprintf("app_%s_%s", os, arch)

		ctr, err := m.buildBinary(ctx, source, goVersion, platform)
		if err != nil {
			return nil, err
		}
		file := ctr.File(filepath.Join("/src", binaryName))
		files = append(files, file)
	}

	return cli.Directory().WithFiles(".", files), nil
}

func (m *GoDagger) buildBinary(ctx context.Context, source *dagger.Directory, goVersion string, platform string) (*Container, error) {
	fmt.Printf("Building binary for %s...\n", platform)

	os := platformFormat.MustParse(string(platform)).OS
	arch := platformFormat.MustParse(string(platform)).Architecture
	binaryName := fmt.Sprintf("app_%s_%s", os, arch)

	cli := dag.Pipeline("go-build-" + binaryName)

	return cli.Container().
		From("golang:"+goVersion).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", m.goModCacheVolume()).
		// run `go mod download` with only go.mod files (re-run only if mod files have changed)
		WithDirectory("/src", source, ContainerWithDirectoryOpts{
			Include: []string{"**/go.mod", "**/go.sum"},
		}).
		WithExec([]string{"go", "mod", "download"}).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", m.goBuildCacheVolume()).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		// run `go build` with all source
		WithMountedDirectory("/src", source).
		WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", binaryName, "."}).
		Sync(ctx)
}

func (m *GoDagger) goModCacheVolume() *CacheVolume {
	return dag.CacheVolume("go-mod")
}

func (m *GoDagger) goBuildCacheVolume() *CacheVolume {
	return dag.CacheVolume("go-build")
}
