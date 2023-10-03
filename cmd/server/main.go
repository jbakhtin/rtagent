package main

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/pkg/closer"
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
	osCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

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

	go func() {
		if err = s.Start(osCtx, cfg); err != nil {
			logger.Info(err.Error())
		}
	}()

	cl, err := closer.New().WithFuncs(s.Shutdown).Build()
	if err != nil {
		logger.Error(err.Error())
	}

	// Gracefully shut down
	<-osCtx.Done()
	withTimeout, cancelShutdownProc := context.WithTimeout(context.Background(), time.Second*cfg.ShutdownTimeout)
	defer cancelShutdownProc()

	err = cl.Close(withTimeout)
	if err != nil {
		logger.Error(err.Error())
	}
}
