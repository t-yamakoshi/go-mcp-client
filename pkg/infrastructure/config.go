package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/t-yamakoshi/go-mcp-client/pkg/config"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/repository"
)

var _ repository.IFConfigRepository = (*ConfigRepositoryImpl)(nil)

// ConfigRepositoryImpl implements the config repository interface
type ConfigRepositoryImpl struct {
	filepath string
}

// NewConfigRepositoryImpl creates a new config repository implementation
func NewConfigRepositoryImpl(filepath string) *ConfigRepositoryImpl {
	return &ConfigRepositoryImpl{
		filepath: filepath,
	}
}

// Load loads configuration from file
func (r *ConfigRepositoryImpl) Load(ctx context.Context) (*config.Config, error) {
	data, err := os.ReadFile(r.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config config.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save saves configuration to file
func (r *ConfigRepositoryImpl) Save(ctx context.Context, config *config.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(r.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
