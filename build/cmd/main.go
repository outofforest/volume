package main

import (
	"github.com/outofforest/build"

	me "github.com/outofforest/volume/build"
)

func main() {
	build.Main("volume2", nil, me.Commands)
}
