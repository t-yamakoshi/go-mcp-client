package main

import (
	"log"
)

func main() {
	cliHandler := InitializeCLIHandler("config.json")
	if err := cliHandler.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
