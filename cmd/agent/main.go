package main

import (
	"github.com/jbakhtin/rtagent/internal/rtagentv2"
	"log"

	"github.com/jbakhtin/rtagent/internal/config"
	"go.uber.org/zap"
)

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

	monitor, err := rtagentv2.NewMonitor(cfg, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	if err := monitor.Start(cfg); err != nil {
		logger.Error(err.Error())
	}
}
