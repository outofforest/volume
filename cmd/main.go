package main

import (
	"context"

	"github.com/outofforest/run"
	"github.com/spf13/pflag"

	"github.com/outofforest/volume"
	"github.com/outofforest/volume/lib/libnet"
)

func main() {
	run.New().Run("volume", func(ctx context.Context) error {
		var address string
		var verbose bool
		pflag.StringVar(&address, "addr", ":8080", "Address to listen on")
		pflag.BoolVarP(&verbose, "verbose", "v", false, "Turn on verbose logging")
		pflag.Parse()

		l, err := libnet.Listen(address)
		if err != nil {
			return err
		}
		defer func() {
			_ = l.Close()
		}()

		return volume.Run(ctx, l)
	})
}
