package main

import (
	"fmt"

	"github.com/jbakhtin/rtagent/internal/config"
	"go.uber.org/zap"

	"github.com/jbakhtin/rtagent/internal/agent"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// TODO: что делать с ошибкой, если не получилось инициализировать логер для логирования ошибок? ;)
		fmt.Println(err)
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
