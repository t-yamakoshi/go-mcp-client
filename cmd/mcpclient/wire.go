//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/repository"
	"github.com/t-yamakoshi/go-mcp-client/pkg/infrastructure"
	"github.com/t-yamakoshi/go-mcp-client/pkg/interfaces"
	"github.com/t-yamakoshi/go-mcp-client/pkg/usecase"
)

func InitializeCLIHandler(configPath string) interfaces.CLIHandler {
	wire.Build(
		infrastructure.NewMCPRepositoryImpl,
		wire.Bind(new(repository.MCPRepository), new(*infrastructure.MCPRepositoryImpl)),
		infrastructure.NewConfigRepositoryImpl,
		wire.Bind(new(repository.ConfigRepository), new(*infrastructure.ConfigRepositoryImpl)),
		usecase.NewMCPUseCase,
		usecase.NewConfigUseCase,
		interfaces.NewCLIHandler,
	)
	return nil
}
