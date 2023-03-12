package main

import (
	"github.com/jbakhtin/rtagent/internal/config"
	"go.uber.org/zap"
	"log"

	"github.com/jbakhtin/rtagent/internal/agent"
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

	// TODO: правильно ли прокидывать логер в структуру? как лучше?
	monitor, err := agent.New(cfg, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	if err := monitor.Start(); err != nil {
		logger.Error(err.Error())
	}
}
