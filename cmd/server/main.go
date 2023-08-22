package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server"
	"go.uber.org/zap"

	_ "net/http/pprof"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg, err := config.NewConfigBuilder().
		WithAllFromFlagsS().
		WithAllFromEnv().
		Build()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	s, err := server.New(cfg)
	if err != nil {
		logger.Error(err.Error())
	}

	ctxServer, cancel := context.WithCancel(context.Background())
	go func() {
		if err = s.Start(ctxServer, cfg); err != nil {
			logger.Info(err.Error())
		}
	}()

	ctxOS, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Gracefully shut down
	<-ctxOS.Done()
	err = s.Shutdown(ctxServer)
	if err != nil {
		logger.Info(err.Error())
	}

	cancel()
	time.Sleep(2 * time.Second)
}
