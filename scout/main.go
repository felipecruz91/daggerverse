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
)

const (
	dockerScoutImage = "index.docker.io/docker/scout-cli:latest"
)

type Scout struct{}

// Cves displays CVEs identified in a container image
// Example usage: `dagger call cves --user $DOCKER_SCOUT_HUB_USER --password=env:DOCKER_SCOUT_HUB_PASSWORD --image alpine:3.18.4 stdout`
func (m *Scout) Cves(
	ctx context.Context,
	user string,
	password *Secret,
	image string) *Container {

	return dag.Container().From(dockerScoutImage).
		WithEnvVariable("DOCKER_SCOUT_HUB_USER", user).
		WithSecretVariable("DOCKER_SCOUT_HUB_PASSWORD", password).
		WithExec([]string{"cves", image})
}
