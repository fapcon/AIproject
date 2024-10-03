package sig

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"studentgit.kata.academy/quant/torque/pkg/logster"
)

var ErrSignalReceived = errors.New("operating system signal")

func ListenSignal(ctx context.Context, logger logster.Logger) error {
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return nil
	case sig := <-sigquit:
		logger.WithField("signal", sig).Infof("Captured signal")
		logger.Infof("Gracefully shutting down server...")
		return ErrSignalReceived
	}
}
