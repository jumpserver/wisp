package common

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	signalCtx        context.Context
	signalCancelFunc context.CancelFunc
	once             sync.Once
)

func GetSignalCtx() context.Context {
	once.Do(func() {
		signalCtx, signalCancelFunc = context.WithCancel(context.Background())
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			signalCancelFunc()
		}()
	})
	return signalCtx
}
