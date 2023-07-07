package build

import (
	"github.com/outofforest/build"
	"github.com/outofforest/buildgo"
)

// Commands is a definition of commands available in build system
var Commands = map[string]build.Command{
	"build": {Fn: buildApp, Description: "Build app"},
	"run":   {Fn: runApp, Description: "Run app"},
}

func init() {
	buildgo.AddCommands(Commands)
}
