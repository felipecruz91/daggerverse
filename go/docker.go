package main

import (
	"context"
	"fmt"
	"runtime"

	platformFormat "github.com/containerd/containerd/platforms"
)

// DockerBuild packages the Go binary into a container image
func (m *GoDagger) DockerBuild(ctx context.Context,
	// bin is the directory containing the cross-platform Go binaries
	// +required
	bin *Directory,
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

	return cli.Container(ContainerOpts{Platform: Platform(platform)}).
		From("alpine:latest").
		WithFile("/bin/app", bin.File(binaryName)).
		WithEntrypoint([]string{"/bin/app"}), nil
}

// DockerPush packages the Go binary into a container image and pushes it to a registry
func (m *GoDagger) DockerPush(ctx context.Context,
	// bin is the directory containing the cross-platform Go binaries
	// +required
	bin *Directory,
	// image is the name of the image to push
	// +required
	image string,
	// platforms is the list of platforms to build the container image for
	// +optional
	// +default=["linux/amd64", "linux/arm64"]
	platforms []string,
	// registryAddress is the address of the container registry
	// +optional
	// +default="docker.io"
	registryAddress string,
	// registryUser is the username for the container registry
	// +required
	registryUser string,
	// registryPassword is the password for the container registry
	// +required
	registryPassword *Secret,
) (string, error) {
	cli := dag.Pipeline("docker-push")

	var platformVariants []*Container
	for _, platform := range platforms {
		ctr, err := m.DockerBuild(ctx, bin, platform)
		if err != nil {
			return "", err
		}
		platformVariants = append(platformVariants, ctr)
	}

	return cli.Container().
		WithRegistryAuth(registryAddress, registryUser, registryPassword).
		Publish(ctx, image, ContainerPublishOpts{PlatformVariants: platformVariants})
}
