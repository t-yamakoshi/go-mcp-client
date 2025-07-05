package provider

import (
	"github.com/google/wire"
	"github.com/t-yamakoshi/go-mcp-client/pkg/usecase"
)

var UsecaseSet = wire.NewSet(
	usecase.NewMCPUsecase,
	usecase.NewConfigUsecase,
)
