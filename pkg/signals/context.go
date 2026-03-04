package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdownContext() (context.Context, <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownSignals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}

	sysChan := make(chan os.Signal, 2)
	closeChan := make(chan struct{}, 1)
	signal.Notify(sysChan, shutdownSignals...)

	go func() {
		<-sysChan
		cancel()
		<-sysChan
		signal.Stop(sysChan)
		closeChan <- struct{}{}
	}()

	return ctx, closeChan
}
