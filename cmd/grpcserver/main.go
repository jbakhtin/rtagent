package main

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/grpcServer"
	"time"
)

func main() {
	timer := time.NewTimer(time.Minute)
	grpcServer := grpcServer.New()

	grpcServer.Run(context.TODO())

	<-timer.C
}