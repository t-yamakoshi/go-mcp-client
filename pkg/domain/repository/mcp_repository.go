package repository

import (
	"context"

	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/response"
)

// MCPRepository defines the interface for MCP operations
type MCPRepository interface {
	// Connection management
	Connect(ctx context.Context, serverURL string) error
	Disconnect() error
	IsConnected() bool

	// Message handling
	SendMessage(ctx context.Context, message *entity.Message) error
	ReceiveMessage(ctx context.Context) (*entity.Message, error)

	// Protocol operations
	Initialize(ctx context.Context, clientInfo entity.ClientInfo) (*response.InitializeResponse, error)
	ListTools(ctx context.Context) ([]entity.Tool, error)
	CallTool(ctx context.Context, toolCall entity.ToolCall) (*entity.ToolResult, error)
}
