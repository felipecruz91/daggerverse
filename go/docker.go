package main

import (
	"context"
	"fmt"
	"runtime"

	platformFormat "github.com/containerd/containerd/platforms"
)

// DockerBuild packages the Go binary into a Docker container
func (m *GoDagger) DockerBuild(ctx context.Context,
	// bin is the directory containing the cross-platform Go binaries
	// +required
	bin *Directory,
	// goVersion is the version of Go to use for building the binary
	// +optional
	// +default="1.22.0"
	goVersion string,
	// platform is the platform to build the binary for
	// +optional
	platform string) (*Container, error) {
	cli := dag.Pipeline("docker-build")

	if platform == "" {
		platform = "linux/" + runtime.GOARCH
	}

	os := platformFormat.MustParse(string(platform)).OS
	arch := platformFormat.MustParse(string(platform)).Architecture
	binaryName := fmt.Sprintf("app_%s_%s", os, arch)

	return cli.Container().From("alpine:latest").
		WithFile("/bin/app", bin.File(binaryName)).
		WithEntrypoint([]string{"/bin/app"}), nil
}
