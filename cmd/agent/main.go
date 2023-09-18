package main

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/agent/aggregator"
	"github.com/jbakhtin/rtagent/internal/agent/jobqueue"
	"github.com/jbakhtin/rtagent/internal/agent/jobsmaker"
	"github.com/jbakhtin/rtagent/internal/agent/messagesender"
	"github.com/jbakhtin/rtagent/internal/agent/sender"
	"github.com/jbakhtin/rtagent/internal/agent/tasker/tasks/limited"
	"github.com/jbakhtin/rtagent/internal/agent/tasker/tasks/once"
	"github.com/jbakhtin/rtagent/internal/agent/tasker/tasks/periodic"
	"github.com/jbakhtin/rtagent/internal/agent/taskmanager"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/pkg/closer"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func init() {
	_, err := fmt.Printf("Build version: %s\nBuild date: %s\nBuild Commit: %s\n", BuildVersion, BuildDate, BuildCommit)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	osCtx, osCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer osCancel()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
		return
	}

	cfg, err := config.NewConfigBuilder().WithAllFromFlagsA().WithAllFromEnv().Build()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	sender, err := sender.New().WithConfig(cfg).Build()
	if err != nil {
		logger.Error(err.Error())
	}

	aggregator, err := aggregator.New().WithDefaultCollectors().WithConfig(cfg).Build()
	if err != nil {
		logger.Error(err.Error())
	}

	task1 := periodic.Task{
		Name: "polling stats",
		Duration: cfg.GetPollInterval(),
		F: aggregator.Pool,
	}

	queue := jobqueue.GetMyQueue()

	jobsMaker := jobsmaker.JobsMaker{
		Jober: queue,
		Slicer: aggregator,
	}

	task2 := periodic.Task{
		Name: "make jobs",
		Duration: cfg.GetReportInterval(),
		F: jobsMaker.Do,
	}

	messageSender := messagesender.MessageSender{
		Sender: sender,
		Jober: queue,
	}

	task3 := limited.Task{
		Name: "send jobs",
		Limit: cfg.RateLimit,
		Duration: time.Second,
		F: messageSender.Do,
	}

	ctx, appCancel := context.WithCancel(osCtx)
	defer appCancel()

	taskManager, err := taskmanager.New().WithFuncs(task1.Do, task2.Do, task3.Do).Build()
	if err != nil {
		logger.Error(err.Error())
	}
	go func () {
		err = taskManager.DoIt(ctx)
		logger.Info(err.Error())
		appCancel()
	}()


	// Gracefully shut down
	select {
	case <-osCtx.Done():
	case <-ctx.Done():
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), time.Second * cfg.GetShutdownTimeout())
	defer cancel()

	task4 := once.Task{
		Name: "send remaining messages",
		F:messageSender.Do,
	}

	closer, err := closer.New().
		WithFuncs(task4.Do).Build()
	if err != nil {
		logger.Info(err.Error())
	}

	err = closer.Close(withTimeout)
	if err != nil {
		logger.Info(err.Error())
	} else {
		logger.Info("shutdown finished successfully")
	}
}
