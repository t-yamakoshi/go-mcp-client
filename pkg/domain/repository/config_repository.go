package repository

import (
	"context"

	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/config"
)

// ConfigRepository defines the interface for configuration operations
type ConfigRepository interface {
	Load(ctx context.Context) (*config.Config, error)
	Save(ctx context.Context, config *config.Config) error
}
