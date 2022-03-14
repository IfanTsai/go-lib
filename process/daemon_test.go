package process_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/docker/docker/pkg/reexec"

	"github.com/IfanTsai/go-lib/process"
)

const testShuttingTime = time.Second * 3

func testDaemon() {
	_ = process.GracefulShutdown(func(ctx context.Context) error {
		<-ctx.Done()

		return nil
	}, testShuttingTime)
}

func TestGracefulShutdown(t *testing.T) {
	t.Parallel()

	reexec.Register("testDaemon", testDaemon)
	require.False(t, reexec.Init())

	cmd := reexec.Command("testDaemon")

	quit := make(chan struct{}, 1)
	require.NoError(t, cmd.Start())

	// wait util subprocess work done
	time.Sleep(time.Second)

	go func() {
		_ = cmd.Wait()
		quit <- struct{}{}
	}()

	start := time.Now()
	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
	<-quit
	end := time.Now()
	require.WithinDuration(t, start.Add(testShuttingTime), end, time.Second)
}
