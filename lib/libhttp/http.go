package libhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/outofforest/logger"
	"github.com/ridge/must"
	"github.com/ridge/parallel"
	"go.uber.org/zap"
)

// Run starts http server
func Run(ctx context.Context, listener net.Listener, handler http.Handler) error {
	return parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		log := logger.Get(ctx)
		errorLog, err := zap.NewStdLogAt(log, zap.WarnLevel)
		must.OK(err)

		server := http.Server{
			Handler:  handler,
			ErrorLog: errorLog,
		}
		spawn("server", parallel.Fail, func(ctx context.Context) error {
			log.Info("Serving requests")
			err := server.Serve(listener)
			if errors.Is(err, http.ErrServerClosed) && ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		})
		spawn("shutdownHandler", parallel.Fail, func(ctx context.Context) error {
			<-ctx.Done()
			log.Info("Shutting down")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// nolint: contextcheck
			err := server.Shutdown(shutdownCtx)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("Shutdown failed", zap.Error(err))
				return err
			}

			log.Info("Shutdown complete")
			return nil
		})
		return nil
	})
}
