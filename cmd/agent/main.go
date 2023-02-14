package main

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agent"
	"time"
)

const (
	pollInterval = time.Second * 2
	reportInterval = time.Second * 10
)

func main() {
	ctx := context.Background()
	monitor := agent.NewMonitor(ctx, pollInterval, reportInterval)

	monitor.Start()
	defer monitor.Stop()

	timer := time.NewTimer(time.Second * 60 * 60)
	<-timer.C
}
