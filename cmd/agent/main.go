package main

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/agent"
	"time"
)

const (
	pollInterval = time.Second * 2
	reportInterval = time.Second * 10
)

func main() {
	err := agent.Start(pollInterval, reportInterval)
	if err != nil {
		fmt.Println(err)
	}
}
