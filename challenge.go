package main

import (
	"context"

	"dagger/brokelabs-ci/internal/dagger"
)

// ChallengeLint validates a challenge package with `brokelabs challenge lint`
// (the same validation the operator runs at install time).
func (m *BrokelabsCi) ChallengeLint(
	ctx context.Context,
	// Challenge repository root (with challenge.yaml).
	// +ignore=[".git", ".idea", ".vscode"]
	source *dagger.Directory,
	// Version of the BrokeLabs CLI to install.
	// +optional
	// +default="v0.3.0"
	cliVersion string,
	// Local checkout of brokelabs-cli built instead of the published version.
	// +optional
	// +ignore=["bin", "dist", ".idea", ".vscode"]
	cli *dagger.Directory,
	// Local checkout of brokelabs-core-go used by the local CLI.
	// +optional
	// +ignore=["bin", ".idea", ".vscode"]
	coreGo *dagger.Directory,
) (string, error) {
	var ctr *dagger.Container
	if cli != nil {
		ctr = m.Go(cli, goVersion, coreGo).Container().
			WithExec([]string{"go", "build", "-o", "/usr/local/bin/brokelabs", "./cmd/brokelabs"})
	} else {
		ctr = m.Go(dag.Directory(), goVersion, nil).Container().
			WithExec([]string{"go", "install", "github.com/Sigeth/brokelabs-cli/cmd/brokelabs@" + cliVersion}).
			WithExec([]string{"cp", "/go/bin/brokelabs", "/usr/local/bin/brokelabs"})
	}
	return ctr.
		WithMountedDirectory("/challenge", source).
		WithExec([]string{"brokelabs", "challenge", "lint", "/challenge"}).
		Stdout(ctx)
}
