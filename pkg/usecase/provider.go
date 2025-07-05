package usecase

import "github.com/google/wire"

var UsecaseSet = wire.NewSet(
	NewMCPUsecase,
	NewConfigUsecase,
)
