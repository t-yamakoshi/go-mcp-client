package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/config"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/repository"
)

// ConfigUseCase defines the interface for configuration business logic
type ConfigUseCase interface {
	LoadConfiguration(ctx context.Context, configPath string) (*config.Config, error)
	SaveConfiguration(ctx context.Context, config *config.Config, configPath string) error
	GetDefaultConfiguration(ctx context.Context) *config.Config
	ValidateConfiguration(ctx context.Context, config *config.Config) error
	UpdateConfiguration(ctx context.Context, config *config.Config, updates map[string]interface{}) error
}

// configUseCase implements the ConfigUseCase interface
type configUseCase struct {
	configRepo repository.ConfigRepository
}

// NewConfigUseCase creates a new configuration use case
func NewConfigUseCase(configRepo repository.ConfigRepository) ConfigUseCase {
	return &configUseCase{
		configRepo: configRepo,
	}
}

// LoadConfiguration loads configuration from file or creates default
func (uc *configUseCase) LoadConfiguration(ctx context.Context, configPath string) (*config.Config, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); err == nil {
		// File exists, load it
		config, err := uc.configRepo.Load(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
		return config, nil
	}

	// File doesn't exist, create default configuration
	defaultConfig := uc.GetDefaultConfiguration(ctx)

	// Save default config for future use
	if err := uc.SaveConfiguration(ctx, defaultConfig, configPath); err != nil {
		return nil, fmt.Errorf("failed to save default config: %w", err)
	}

	return defaultConfig, nil
}

// SaveConfiguration saves configuration to file
func (uc *configUseCase) SaveConfiguration(ctx context.Context, config *config.Config, configPath string) error {
	// Validate configuration before saving
	if err := uc.ValidateConfiguration(ctx, config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return uc.configRepo.Save(ctx, config)
}

// GetDefaultConfiguration returns the default configuration
func (uc *configUseCase) GetDefaultConfiguration(ctx context.Context) *config.Config {
	return &config.Config{
		ServerURL: "ws://localhost:3000",
		ClientInfo: entity.ClientInfo{
			Name:    "go-mcp-client",
			Version: "1.0.0",
		},
		LogLevel: "info",
	}
}

// ValidateConfiguration validates the configuration
func (uc *configUseCase) ValidateConfiguration(ctx context.Context, config *config.Config) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if config.ServerURL == "" {
		return fmt.Errorf("server URL cannot be empty")
	}

	if config.ClientInfo.Name == "" {
		return fmt.Errorf("client name cannot be empty")
	}

	if config.ClientInfo.Version == "" {
		return fmt.Errorf("client version cannot be empty")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug":   true,
		"info":    true,
		"warning": true,
		"error":   true,
	}

	if !validLogLevels[config.LogLevel] {
		return fmt.Errorf("invalid log level: %s", config.LogLevel)
	}

	return nil
}

// UpdateConfiguration updates specific configuration fields
func (uc *configUseCase) UpdateConfiguration(ctx context.Context, config *config.Config, updates map[string]interface{}) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Apply updates
	if serverURL, ok := updates["server_url"].(string); ok {
		config.ServerURL = serverURL
	}

	if clientName, ok := updates["client_name"].(string); ok {
		config.ClientInfo.Name = clientName
	}

	if clientVersion, ok := updates["client_version"].(string); ok {
		config.ClientInfo.Version = clientVersion
	}

	if logLevel, ok := updates["log_level"].(string); ok {
		config.LogLevel = logLevel
	}

	// Validate the updated configuration
	return uc.ValidateConfiguration(ctx, config)
}
