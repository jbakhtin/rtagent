package main

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/servers"
)

func main() {
	error := servers.Start()
	if error != nil {
		fmt.Println(error)
	}
}
