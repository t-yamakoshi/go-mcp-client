package infrastructure

import "github.com/google/wire"

var InfrastructureSet = wire.NewSet(
	NewMCPRepositoryImpl,
	NewConfigRepositoryImpl,
)
