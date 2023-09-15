package main

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/agent"
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

	monitor, err := agent.New(cfg, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	if err := monitor.Start(cfg); err != nil {
		logger.Error(err.Error())
	}
}
