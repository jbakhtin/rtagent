package main

import (
	"fmt"

	"github.com/jbakhtin/rtagent/internal/server"
)

const serverDomain = "127.0.0.1"
const serverPort = "8080"

func main() {
	serverAddress := fmt.Sprintf("%s:%s", serverDomain, serverPort)

	s, err := server.New(serverAddress)
	if err != nil {
		fmt.Println(err)
	}

	err = s.Start()
	if err != nil {
		fmt.Println(err) // TODO: реализовать логирование ошибок
	}
}
