package process

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func GracefulShutdown(shutdown func(ctx context.Context) error, shuttingTime time.Duration) error {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), shuttingTime)
	defer cancel()

	if err := shutdown(ctx); err != nil {
		return errors.Wrap(err, "force shutdown")
	}

	return nil
}
