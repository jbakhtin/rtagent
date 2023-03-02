package main

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server"
)

const serverDomain = "127.0.0.1"
const serverPort = "8080"

func main() {
	cfg, err := config.NewConfigBuilder().
		WithAddressFromFlag().
		WithRestoreFromFlag().
		WithStoreFileFromFlag().
		WithAllFromEnv().
		Build()

	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: добавить логер для сервера

	s, err := server.New(cfg)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err = s.Start(ctx, cfg); err != nil {
		fmt.Println(err)
		cancel()
	}
}
