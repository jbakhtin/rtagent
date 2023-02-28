package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/jbakhtin/rtagent/internal/config"
	"log"

	"github.com/jbakhtin/rtagent/internal/server"
)

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	// TODO: добавить логер для сервера

	s, err := server.New(cfg)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err = s.Start(ctx, cfg); err != nil {
		log.Println(err)
		cancel()
	}
	log.Println("сервер остановлен")
	cancel()
}
