package main

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/jbakhtin/rtagent/internal/agent"
)

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
	serverDomain   = "http://127.0.0.1"
	serverPort     = "8080"
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
