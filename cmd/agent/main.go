package main

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/agent/aggregator"
	"github.com/jbakhtin/rtagent/internal/agent/jobqueue"
	"github.com/jbakhtin/rtagent/internal/agent/sender"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"

	"github.com/jbakhtin/rtagent/pkg/closer"
	"github.com/jbakhtin/rtagent/pkg/ratelimiter"
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

	cfg, err := config.NewConfigBuilder().WithAllFromFlagsA().WithAllFromEnv().Build()
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
	go func() {
		ticker := time.NewTicker(cfg.GetPollInterval())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				go aggregator.Pool(osCtx)
			case err := <-aggregator.Err():
				errorChan<- err
			case <-osCtx.Done():
			}
		}
	}()

	queue := jobqueue.GetMyQueue()
	poolingStats := func(ctx context.Context) error {
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
					node := jobqueue.GetQNode(key, metric)
					queue.Enqueue(node)
				}

			case <-ctx.Done():
				return errors.Wrap(ctx.Err(), "pooling stats")
			}
		}
	}

	go func() {
		err = poolingStats(osCtx)
		if err != nil {
			logger.Info(err.Error())
		}
	}()


	sendMessages := func(ctx context.Context) error {
		ticker := time.NewTicker(cfg.GetReportInterval())
		defer ticker.Stop()

		rl := ratelimiter.New(time.Second, cfg.RateLimit)
		rl.Run()

		for {
			select {
			case <-ctx.Done():
				return errors.Wrap(ctx.Err(), "send messages")
			case <-rl.Wait():
				if queue.IsEmpty() {
					continue
				}

				node := queue.Dequeue()
				go func(key string, metric types.Metricer) {
					err = sender.Send(key, metric)
					if err != nil {
						errorChan<- err
					}
				}(node.Key(), node.Value())
			}
		}
	}

	go func() {
		err = sendMessages(osCtx)
		if err != nil {
			logger.Info(err.Error())
		}
	}()

	// check error count
	var errCount int
	checkErrors := func(ctx context.Context) error {
		for {
			select {
			case err = <-errorChan:
				errCount++
				logger.Info(err.Error())

				if errCount > cfg.GetAcceptableCountAgentErrors() {
					logger.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", cfg.GetAcceptableCountAgentErrors()))
					cancel()
				}
			case <-ctx.Done():
				return errors.Wrap(ctx.Err(), "check error count")
			}
		}
	}

	go func() {
		err = checkErrors(osCtx)
		if err != nil {
			logger.Info(err.Error())
		}
	}()

	// Gracefully shut down
	<-osCtx.Done()
	withTimeout, cancel := context.WithTimeout(context.Background(), time.Second * cfg.GetShutdownTimeout())
	defer cancel()

	closer, err := closer.New().
		WithFuncs(sendMessages).Build()
	if err != nil {
		logger.Info(err.Error())
	}

	err = closer.Close(withTimeout)
	if err != nil {
		logger.Info(err.Error())
	}
}
