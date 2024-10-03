package logster

import (
	"context"
	"io"
	"sync"
	"time"
)

func Elapsed(logger Logger, startedAt time.Time, msg string, args ...interface{}) {
	logger.WithPrefix(LibPrefix).WithField("elapsed_ms", elapsedMs(startedAt)).Infof(msg, args...)
}

func FatalIfError(logger Logger, err error, msg string, args ...interface{}) {
	if err != nil {
		logger.WithError(err).Fatalf(msg, args...)
	}
}

func LogIfError(logger Logger, err error, msg string, args ...interface{}) error {
	if err != nil {
		logger.WithError(err).Errorf(msg, args...)
	}
	return err
}

func InitResource(logger Logger, resource interface{}, err error) {
	logger.Debugf("Init %T", resource)
	if err != nil {
		logger.WithError(err).Fatalf("Can't init %T", resource)
	}
}

func CloseResource(logger Logger, resource io.Closer) {
	logger.Debugf("Close %T", resource)
	if err := resource.Close(); err != nil {
		logger.WithError(err).Errorf("Can't close %T", resource)
	}
}

func DurationToMs(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / 1000.0 / 1000.0
}

func elapsedMs(since time.Time) float64 {
	return DurationToMs(time.Since(since))
}

func LogShutdownDuration(ctx context.Context, logger Logger) func() {
	var mu sync.Mutex
	var shutdownTime time.Time

	go func() {
		<-ctx.Done()
		mu.Lock()
		shutdownTime = time.Now()
		mu.Unlock()
	}()

	return func() {
		mu.Lock()
		Elapsed(logger, shutdownTime, "Shutdown duration")
		mu.Unlock()
	}
}
