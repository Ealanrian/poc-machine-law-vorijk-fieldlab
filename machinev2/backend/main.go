package main

import (
	"log"

	"github.com/Ealanrian/poc-machine-law/machinev2/backend/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
