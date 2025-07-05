package main

import (
	"log"

	"github.com/t-yamakoshi/go-mcp-client/cmd/mcpclient/di"
)

func main() {
	cliHandler := di.InitializeCLIHandler("config.json")
	if err := cliHandler.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
