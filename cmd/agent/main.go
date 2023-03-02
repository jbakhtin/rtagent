package main

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"time"

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
	cfg, err := config.NewConfigBuilder().
		WithAddressFromFlag().
		WithPollIntervalFromFlag().
		WithReportIntervalFromFlag().
		WithAllFromEnv().
		Build()

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
