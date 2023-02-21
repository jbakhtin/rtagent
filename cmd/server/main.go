package main

import (
	"fmt"

	"github.com/jbakhtin/rtagent/internal/server"
)

const serverDomain = "127.0.0.1"
const serverPort = "8080"

func main() {
	serverAddress := fmt.Sprintf("%s:%s", serverDomain, serverPort)

	err := server.Start(serverAddress)
	if err != nil {
		fmt.Println(err) // TODO: реализовать логирование ошибок
	}
}
