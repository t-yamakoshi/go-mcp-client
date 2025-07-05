package service

import (
	"context"

	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/config"
)

// ConfigService defines the interface for configuration business logic
type ConfigService interface {
	LoadConfiguration(ctx context.Context, filepath string) (*config.Config, error)
	SaveConfiguration(ctx context.Context, config *config.Config, filepath string) error
	GetDefaultConfiguration(ctx context.Context) *config.Config
	ValidateConfiguration(ctx context.Context, config *config.Config) error
}
