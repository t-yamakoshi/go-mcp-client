package domain

import (
	"context"

	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
)

// MCPService defines the interface for MCP business logic
type MCPService interface {
	// Connection management
	EstablishConnection(ctx context.Context, serverURL string) error
	CloseConnection(ctx context.Context) error
	GetConnectionStatus(ctx context.Context) entity.ConnectionStatus

	// Protocol operations
	InitializeProtocol(ctx context.Context, clientInfo entity.ClientInfo) (*InitializeResponse, error)
	GetAvailableTools(ctx context.Context) ([]entity.Tool, error)
	ExecuteTool(ctx context.Context, toolCall entity.ToolCall) (*entity.ToolResult, error)

	// Message handling
	HandleIncomingMessage(ctx context.Context, message *entity.Message) error
	SendOutgoingMessage(ctx context.Context, message *entity.Message) error
}

// ConfigService defines the interface for configuration business logic
type ConfigService interface {
	LoadConfiguration(ctx context.Context, filepath string) (*Config, error)
	SaveConfiguration(ctx context.Context, config *Config, filepath string) error
	GetDefaultConfiguration(ctx context.Context) *Config
	ValidateConfiguration(ctx context.Context, config *Config) error
}
