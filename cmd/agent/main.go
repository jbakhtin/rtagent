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
	err := agent.Start(serverDomain+":"+serverPort, pollInterval, reportInterval)
	if err != nil {
		fmt.Println(err)
	}
}
