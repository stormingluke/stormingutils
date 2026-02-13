// Stormingutils builds a multi-tool CLI container with turso, flyctl, jq, and curl.
// Used by the stormingplatform pipeline for all CLI operations.

package main

import (
	"context"
	"fmt"

	"dagger/stormingutils/internal/dagger"
)

const (
	BaseImage = "debian:bookworm-slim"
	GhcrRepo  = "ghcr.io/stormingluke/stormingutils"
)

type Stormingutils struct{}

// Build creates the stormingutils container with turso, flyctl, jq, and curl installed.
func (m *Stormingutils) Build() *dagger.Container {
	return dag.Container(dagger.ContainerOpts{Platform: "linux/amd64"}).
		From(BaseImage).
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{
			"apt-get", "install", "-y", "--no-install-recommends",
			"ca-certificates", "curl", "jq", "unzip", "xz-utils",
		}).
		WithExec([]string{"rm", "-rf", "/var/lib/apt/lists/*"}).
		WithExec([]string{"sh", "-c", "curl -fsSL https://fly.io/install.sh | FLYCTL_INSTALL=/usr/local sh"}).
		WithExec([]string{"sh", "-c", "curl -sSfL https://get.tur.so/install.sh | bash && mv /root/.turso/turso /usr/local/bin/turso && rm -rf /root/.turso"}).
		WithExec([]string{"flyctl", "version"}).
		WithExec([]string{"turso", "--version"}).
		WithEntrypoint([]string{"/bin/bash"})
}

// Publish builds and pushes the stormingutils container to GHCR.
func (m *Stormingutils) Publish(
	ctx context.Context,
	// GitHub token with packages:write scope
	ghToken *dagger.Secret,
	// Image tag
	// +optional
	// +default="latest"
	tag string,
) (string, error) {
	ref := fmt.Sprintf("%s:%s", GhcrRepo, tag)
	return m.Build().
		WithRegistryAuth("ghcr.io", "stormingluke", ghToken).
		Publish(ctx, ref)
}
