package main

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/agent"
	"github.com/jbakhtin/rtagent/internal/agent/aggregator"
	"github.com/jbakhtin/rtagent/internal/agent/sender"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/pkg/closer"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func init() {
	_, err := fmt.Printf("Build version: %s\nBuild date: %s\nBuild Commit: %s\n", BuildVersion, BuildDate, BuildCommit)
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
		return
	}

	cfg, err := config.NewConfigBuilder().
		WithAllFromFlagsA().
		WithAllFromEnv().
		Build()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	sender, err := sender.New().WithConfig(cfg).Build()
	if err != nil {
		logger.Error(err.Error())
	}

	aggregator, err := aggregator.New().WithDefaultCollectors().WithConfig(cfg).Build()
	if err != nil {
		logger.Error(err.Error())
	}

	agent, err := agent.New().WithConfig(cfg).WithSender(sender).WithAggregator(aggregator).WithSoftShuttingDown().Build()
	if err != nil {
		logger.Error(err.Error())
	}

	agent.Run(osCtx, cfg)

	// check error count
	var errCount int
	go func() {
		for {
			select {
			case err = <-agent.Err():
				errCount++
				logger.Info(err.Error())

				if errCount > cfg.GetAcceptableCountAgentErrors() {
					logger.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", cfg.GetAcceptableCountAgentErrors()))
					cancel()
				}
			case <-osCtx.Done():
				logger.Info("завершаем работу агента")
				return
			}
		}
	}()

	// Gracefully shut down
	<-osCtx.Done()
	withTimeout, cancel := context.WithTimeout(context.Background(), time.Second * cfg.GetShutdownTimeout())
	defer cancel()

	closer, err := closer.New().
		WithFuncs(
			agent.Close,
		).Build()
	if err != nil {
		logger.Info(err.Error())
	}

	err = closer.Close(withTimeout)
	if err != nil {
		logger.Info(err.Error())
	}
}
