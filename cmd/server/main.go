package main

import (
	"context"
	"fmt"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		return
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

	ctx, cancel := context.WithCancel(context.Background())
	if err = s.Start(ctx, cfg); err != nil {
		logger.Error(err.Error())
		cancel() // TODO: реализовать мягкое завершение всех процессов
	}
}
