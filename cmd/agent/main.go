package main

import (
	"fmt"
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
	serverAddress := fmt.Sprintf("%s:%s", serverDomain, serverPort)

	monitor, err := agent.New(serverAddress, pollInterval, reportInterval)
	if err != nil {
		fmt.Println(err)
	}

	if err := monitor.Start(); err != nil {
		fmt.Println(err)
		fmt.Println("Test")
	}
}
