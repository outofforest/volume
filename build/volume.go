package build

import (
	"context"
	"os/exec"

	"github.com/outofforest/build"
	"github.com/outofforest/buildgo"
	"github.com/outofforest/libexec"
)

func buildApp(ctx context.Context) error {
	return buildgo.GoBuildPkg(ctx, "cmd", "bin/volume-app", true)
}

func runApp(ctx context.Context, deps build.DepsFunc) error {
	deps(buildApp)
	return libexec.Exec(ctx, exec.Command("./bin/volume-app"))
}
