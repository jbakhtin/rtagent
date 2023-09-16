package main

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/agent"
	aggregator2 "github.com/jbakhtin/rtagent/internal/agent/aggregator"
	"github.com/jbakhtin/rtagent/internal/agent/sender"
	"log"

	"github.com/jbakhtin/rtagent/internal/config"
	"go.uber.org/zap"
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

	aggregator, err := aggregator2.New().WithDefaultCollectors().WithConfig(cfg).Build()
	if err != nil {
		logger.Error(err.Error())
	}

	agent, err := agent.New().WithConfig(cfg).WithSender(sender).WithAggregator(aggregator).Build()
	if err != nil {
		logger.Error(err.Error())
	}

	if err := agent.Start(cfg); err != nil {
		logger.Error(err.Error())
	}
}
