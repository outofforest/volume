package main

import (
	"context"

	"github.com/outofforest/logger"
	"github.com/outofforest/run"
)

func main() {
	run.Service("volume", nil, func(ctx context.Context) error {
		log := logger.Get(ctx)
		log.Info("test")
		return nil
	})
}
