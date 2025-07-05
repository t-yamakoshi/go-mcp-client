//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/t-yamakoshi/go-mcp-client/cmd/mcpclient/di/provider"
)

func InitializeCLIHandler(configPath string) *provider.CliHandler {
	wire.Build(
		provider.InfrastructureSet,
		provider.UsecaseSet,
		provider.MessageSet,
		provider.CLISet,
	)
	return nil
}
