package main

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/rtagent"
	"time"
)

func main() {
	ctx := context.Background()

	monitor := rtagent.NewMonitor(ctx, time.Second * 2, time.Second * 10)

	monitor.Start()
	defer monitor.Stop()

	timer := time.NewTimer(time.Second * 60 * 60)
	<-timer.C
}
