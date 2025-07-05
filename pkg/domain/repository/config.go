package repository

import (
	"context"

	"github.com/t-yamakoshi/go-mcp-client/pkg/config"
)

// IFConfigRepository defines the interface for configuration operations
type IFConfigRepository interface {
	Load(ctx context.Context) (*config.Config, error)
	Save(ctx context.Context, config *config.Config) error
}
