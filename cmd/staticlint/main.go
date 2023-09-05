package main

import (
	"github.com/jbakhtin/rtagent/pkg/multichecker"
	"log"
)

func main() {
	err := multichecker.Start()
	if err != nil {
		log.Fatal(err)
	}
}
