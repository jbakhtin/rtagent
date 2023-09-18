package main

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/agent/aggregator"
	"github.com/jbakhtin/rtagent/internal/agent/sender"
	"github.com/jbakhtin/rtagent/internal/agentv2/workerPool"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"
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

	errorChan := make(chan error)
	wp, _ := workerPool.New()
	// Task 1
	wp.AddJob(func() error {
		ticker := time.NewTicker(cfg.GetPollInterval())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				go aggregator.Pool(osCtx)
			case err := <-aggregator.Err():
				errorChan<- err
			case <-osCtx.Done():
				return nil
			}
		}
	})

	// Task 1
	wp.AddJob(func() error {
		ticker := time.NewTicker(cfg.GetReportInterval())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats, err := aggregator.GetAll()
				if err != nil {
					errorChan<- err
				}

				for key, metric := range stats {
					go func(key string, metric types.Metricer) {
						err = sender.Send(key, metric)
						if err != nil {
							errorChan<- err
						}
					}(key, metric)
				}

			case <-osCtx.Done():
				return nil
			}
		}
	})

	// check error count
	var errCount int
	go func() {
		for {
			select {
			case err = <-errorChan:
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
		WithFuncs().Build()
	if err != nil {
		logger.Info(err.Error())
	}

	err = closer.Close(withTimeout)
	if err != nil {
		logger.Info(err.Error())
	}
}
