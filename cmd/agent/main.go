package main

import (
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/jbakhtin/rtagent/internal/agent"
)

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
	serverDomain   = "http://127.0.0.1"
	serverPort     = "8080"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		return
	}

	serverAddress := fmt.Sprintf("%s:%s", serverDomain, serverPort)
	monitor, err := agent.New(serverAddress, pollInterval, reportInterval, logger)
	if err != nil {
		fmt.Println(err)
	}

	if err := monitor.Start(); err != nil {
		logger.Error(err.Error())
	}
}
