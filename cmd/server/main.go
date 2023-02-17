package main

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/server"
)

func main() {
	err := server.Start() // TODO: нужно ли пробрасывать контекст, может ли это пригодится?
	if err != nil {
		fmt.Println(err) // TODO: реализовать логирование ошибок
	}
}
