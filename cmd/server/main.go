package main

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/grpcServer"
	"github.com/jbakhtin/rtagent/internal/storage"
	"github.com/jbakhtin/rtagent/internal/storage/filestorage"
	"github.com/jbakhtin/rtagent/pkg/closer"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server"
	"go.uber.org/zap"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"

	buildFmt = "Build version: %s\nBuild date: %s\nBuild commit: %s\n"
)

var (
	repository storage.MetricRepository
	logger *zap.Logger
	cfg *config.Config
	grpc *grpcServer.Server
	http server.MainServer
	clr *closer.Closer
)

func init() {
	var err error

	if _, err = fmt.Printf(buildFmt, BuildVersion, BuildDate, BuildCommit); err != nil {
		log.Fatal(errors.Wrap(err, "set app configuration info"))
	}

	if logger, err = zap.NewDevelopment(); err != nil {
		log.Fatal(errors.Wrap(err, "init logger"))
	}

	if cfg, err = config.NewConfigBuilder().WithAllFromFlagsS().WithAllFromEnv().Build(); err != nil {
		log.Fatal(errors.Wrap(err, "init config"))
	}

	if repository, err = storage.New().File(cfg).Build(); err != nil {
		log.Fatal(errors.Wrap(err, "init repository"))
	}

	if grpc = grpcServer.New(*cfg, repository); err != nil {
		log.Fatal(errors.Wrap(err, "init grpc server"))
	}

	if http, err = server.New(*cfg, repository); err != nil {
		log.Fatal(errors.Wrap(err, "init http server"))
	}

	if clr, err = closer.New().WithFuncs(http.Shutdown).Build(); err != nil {
		log.Fatal(errors.Wrap(err, "init closer"))
	}
}

func main() {
	osCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	switch temp := repository.(type) {
	case *filestorage.FileStorage:
		// ToDo: need to move bckuper to Facade
		if err := temp.Start(osCtx, *cfg); err != nil {
			log.Fatal(errors.Wrap(err, "start backup storage"))
		}
	// ...
	}

	if err := grpc.Run(osCtx); err != nil {
		log.Fatal(errors.Wrap(err, "start grpc server"))
	}

	if err := http.Start(osCtx, *cfg); err != nil {
		log.Fatal(errors.Wrap(err, "start http server"))
	}

	// Gracefully shut down
	<-osCtx.Done()
	withTimeout, cancelShutdownProc := context.WithTimeout(context.Background(), time.Second*cfg.ShutdownTimeout)
	defer cancelShutdownProc()

	if err := clr.Close(withTimeout); err != nil {
		log.Fatal(errors.Wrap(err, "shutdown"))
	}
}
