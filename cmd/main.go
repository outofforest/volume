package main

import (
	"context"

	"github.com/outofforest/logger"
	"github.com/outofforest/run"
	"github.com/outofforest/volume"
	"github.com/outofforest/volume/lib/libnet"
	"github.com/spf13/pflag"
)

func main() {
	var address string
	var verbose bool
	pflag.StringVar(&address, "addr", ":8080", "Address to listen on")
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Turn on verbose logging")
	pflag.Parse()

	if !verbose {
		logger.VerboseOff()
	}

	run.Service("volume", nil, func(ctx context.Context) error {
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
