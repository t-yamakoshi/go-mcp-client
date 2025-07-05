package provider

import (
	"github.com/google/wire"
	"github.com/t-yamakoshi/go-mcp-client/pkg/infrastructure"
)

var InfrastructureSet = wire.NewSet(
	infrastructure.NewMCPRepositoryImpl,
	infrastructure.NewConfigRepositoryImpl,
)
