package build

import (
	"github.com/outofforest/build"
	"github.com/outofforest/buildgo"
)

// Commands is a definition of commands available in build system
var Commands = map[string]build.Command{
	"build": build.Command{Fn: buildApp, Description: "Build app"},
	"run":   build.Command{Fn: runApp, Description: "Run app"},
}

func init() {
	buildgo.AddCommands(Commands)
}
