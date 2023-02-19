package main

import (
	"fmt"

	"github.com/jbakhtin/rtagent/internal/server"
)

const serverAddress = "127.0.0.1"
const serverPort = "8080"

func main() {
	err := server.Start(serverAddress + ":" + serverPort) // TODO: нужно ли пробрасывать контекст, может ли это пригодится?
	if err != nil {
		fmt.Println(err) // TODO: реализовать логирование ошибок
	}
}
