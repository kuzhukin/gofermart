package main

import (
	"fmt"
	"gophermart/internal/apiserver"
	"gophermart/internal/zlog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serverStopTimeout = time.Second * 30

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	defer func() {
		_ = zlog.Logger.Sync()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	srvr, err := apiserver.StartNew()
	if err != nil {
		return fmt.Errorf("start gophermart server, err=%w", err)
	}

	sig := <-sigs
	zlog.Logger.Infof("Stop server by osSignal=%v", sig)
	if err := srvr.Stop(); err != nil {
		return fmt.Errorf("stop server, err=%w", err)
	}

	select {
	case <-srvr.Wait():
		zlog.Logger.Infof("Server stopped")
	case <-time.After(serverStopTimeout):
		zlog.Logger.Infof("Server stopped by timeout=%v", serverStopTimeout)
	}

	return nil
}
