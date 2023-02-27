package main

import (
	"fmt"

	"github.com/jbakhtin/rtagent/internal/server"
)

const serverDomain = "127.0.0.1"
const serverPort = "8080"

func main() {
	// TODO: добавить логер для сервера
	serverAddress := fmt.Sprintf("%s:%s", serverDomain, serverPort)

	s, err := server.New(serverAddress)
	if err != nil {
		fmt.Println(err)
	}

	if err = s.Start(); err != nil {
		fmt.Println(err)
	}
}
