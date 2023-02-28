package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/jbakhtin/rtagent/internal/config"
	"go.uber.org/zap"
	"log"

	"github.com/jbakhtin/rtagent/internal/agent"
)

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		return
	}

	monitor, err := agent.New(cfg, logger)
	if err != nil {
		fmt.Println(err)
	}

	if err := monitor.Start(); err != nil {
		logger.Error(err.Error())
	}
}
