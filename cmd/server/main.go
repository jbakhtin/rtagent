package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/jbakhtin/rtagent/internal/config"
	"log"

	"github.com/jbakhtin/rtagent/internal/server"
)

const serverDomain = "127.0.0.1"
const serverPort = "8080"

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

	if err = s.Start(); err != nil {
		fmt.Println(err)
	}
}
