package main

import (
	"log"

	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/backend/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
