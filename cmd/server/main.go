package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server"
	"go.uber.org/zap"

	_ "net/http/pprof"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func init() {
	_, err := fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", BuildVersion, BuildDate, BuildCommit)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
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
