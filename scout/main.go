// A generated module for Scout functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/scout/internal/dagger"
	"errors"
	"fmt"
)

const (
	dockerScoutImage = "index.docker.io/docker/scout-cli:latest"

	// vulnExitCode is the exit code returned by Docker Scout when vulnerabilities are found
	vulnExitCode = 2
)

type Scout struct{}

// DockerScoutCves displays CVEs identified in a container image
func (m *Scout) DockerScoutCves(
	ctx context.Context,
	// source is the directory containing the image tarball
	// +required
	source *Directory,
	// dockerScoutHubUser is the username for Docker Scout Hub
	// +required
	dockerScoutHubUser string,
	// dockerScoutHubPassword is the password for Docker Scout Hub
	// +required
	dockerScoutHubPassword *Secret,
	// tarballPath is the path to the tarball containing the container image
	// +required
	tarballPath string,
	// onlySeverity is the severity of vulnerabilities to filter by
	// +optional
	// +default="critical,high"
	onlySeverity string,
	// exitCode returns '2' if vulnerabilities are found
	// +optional
	// +default=false
	exitCode bool,
) (*Container, error) {
	cli := dag.Pipeline("docker-scout-cves")

	args := []string{"cves", "archive://" + tarballPath}
	if onlySeverity != "" {
		args = append(args, "--only-severity", onlySeverity)
	}
	if exitCode {
		args = append(args, "--exit-code")
	}

	ctr, err := cli.Container().From(dockerScoutImage).
		WithEnvVariable("DOCKER_SCOUT_HUB_USER", dockerScoutHubUser).
		WithSecretVariable("DOCKER_SCOUT_HUB_PASSWORD", dockerScoutHubPassword).
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec(args).
		Sync(ctx)

	var e *dagger.ExecError
	if exitCode && errors.As(err, &e) {
		if e.ExitCode == vulnExitCode {
			return ctr, fmt.Errorf("failing function because --exit-code was provided and vulnerabilities were found in container image, see stdout above for details")
		}
		return ctr, err
	}

	return ctr, nil
}
