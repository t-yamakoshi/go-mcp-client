package service

import (
	"context"

	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/response"
)

// MCPService defines the interface for MCP business logic
type MCPService interface {
	// Connection management
	EstablishConnection(ctx context.Context, serverURL string) error
	CloseConnection(ctx context.Context) error
	GetConnectionStatus(ctx context.Context) entity.ConnectionStatus

	// Protocol operations
	InitializeProtocol(ctx context.Context, clientInfo entity.ClientInfo) (*response.InitializeResponse, error)
	GetAvailableTools(ctx context.Context) ([]entity.Tool, error)
	ExecuteTool(ctx context.Context, toolCall entity.ToolCall) (*entity.ToolResult, error)

	// Message handling
	HandleIncomingMessage(ctx context.Context, message *entity.Message) error
	SendOutgoingMessage(ctx context.Context, message *entity.Message) error
}
