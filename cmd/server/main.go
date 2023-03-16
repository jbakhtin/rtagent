package main

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"
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
	done := make(chan bool, 1)
	// Gracefully shut down
	go func() {
		<-ctxOS.Done()
		err := s.Shutdown(ctxServer)
		if err != nil {
			logger.Info(err.Error())
		}

		cancel()
		time.Sleep(2 * time.Second)

		close(done)
	}()

	<-done
}
