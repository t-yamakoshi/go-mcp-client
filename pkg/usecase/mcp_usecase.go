package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/repository"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/response"
)

// MCPUseCase defines the interface for MCP business logic
// (mockgen用コメント)
//
//go:generate mockgen -source=mcp_usecase.go -destination=mock_mcp_usecase.go -package=usecase
type MCPUseCase interface {
	EstablishConnection(ctx context.Context, serverURL string) error
	CloseConnection(ctx context.Context) error
	GetConnectionStatus(ctx context.Context) entity.ConnectionStatus
	InitializeProtocol(ctx context.Context, clientInfo entity.ClientInfo) (*response.InitializeResponse, error)
	GetAvailableTools(ctx context.Context) ([]entity.Tool, error)
	ExecuteTool(ctx context.Context, toolCall entity.ToolCall) (*entity.ToolResult, error)
	HandleIncomingMessage(ctx context.Context, message *entity.Message) error
	SendOutgoingMessage(ctx context.Context, message *entity.Message) error
	RegisterHandler(method string, handler MessageHandler)
}

// mcpUseCase implements the MCPUseCase interface
type mcpUseCase struct {
	mcpRepo    repository.MCPRepository
	configRepo repository.ConfigRepository
	mu         sync.RWMutex
	handlers   map[string]MessageHandler
	connection *entity.Connection
}

// MessageHandler is a function type for handling incoming messages
type MessageHandler func(*entity.Message) error

// NewMCPUseCase creates a new MCP use case
func NewMCPUseCase(mcpRepo repository.MCPRepository, configRepo repository.ConfigRepository) MCPUseCase {
	return &mcpUseCase{
		mcpRepo:    mcpRepo,
		configRepo: configRepo,
		handlers:   make(map[string]MessageHandler),
		connection: &entity.Connection{
			ID:        uuid.New().String(),
			Status:    entity.ConnectionStatusDisconnected,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// EstablishConnection establishes a connection to the MCP server
func (uc *mcpUseCase) EstablishConnection(ctx context.Context, serverURL string) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	// Update connection status
	uc.connection.Status = entity.ConnectionStatusConnecting
	uc.connection.ServerURL = serverURL
	uc.connection.UpdatedAt = time.Now()

	// Connect to the server
	if err := uc.mcpRepo.Connect(ctx, serverURL); err != nil {
		uc.connection.Status = entity.ConnectionStatusError
		uc.connection.UpdatedAt = time.Now()
		return fmt.Errorf("failed to establish connection: %w", err)
	}

	// Update connection status
	uc.connection.Status = entity.ConnectionStatusConnected
	uc.connection.UpdatedAt = time.Now()

	log.Printf("Successfully connected to MCP server: %s", serverURL)
	return nil
}

// CloseConnection closes the connection to the MCP server
func (uc *mcpUseCase) CloseConnection(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.mcpRepo.Disconnect(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	uc.connection.Status = entity.ConnectionStatusDisconnected
	uc.connection.UpdatedAt = time.Now()

	log.Println("Connection closed successfully")
	return nil
}

// GetConnectionStatus returns the current connection status
func (uc *mcpUseCase) GetConnectionStatus(ctx context.Context) entity.ConnectionStatus {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return uc.connection.Status
}

// InitializeProtocol initializes the MCP protocol
func (uc *mcpUseCase) InitializeProtocol(ctx context.Context, clientInfo entity.ClientInfo) (*response.InitializeResponse, error) {
	uc.mu.RLock()
	if uc.connection.Status != entity.ConnectionStatusConnected {
		uc.mu.RUnlock()
		return nil, fmt.Errorf("not connected to server")
	}
	uc.mu.RUnlock()

	response, err := uc.mcpRepo.Initialize(ctx, clientInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize protocol: %w", err)
	}

	log.Printf("Protocol initialized with server: %s v%s",
		response.ServerInfo.Name, response.ServerInfo.Version)
	return response, nil
}

// GetAvailableTools retrieves available tools from the server
func (uc *mcpUseCase) GetAvailableTools(ctx context.Context) ([]entity.Tool, error) {
	uc.mu.RLock()
	if uc.connection.Status != entity.ConnectionStatusConnected {
		uc.mu.RUnlock()
		return nil, fmt.Errorf("not connected to server")
	}
	uc.mu.RUnlock()

	tools, err := uc.mcpRepo.ListTools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tools: %w", err)
	}

	log.Printf("Retrieved %d available tools", len(tools))
	return tools, nil
}

// ExecuteTool executes a tool on the server
func (uc *mcpUseCase) ExecuteTool(ctx context.Context, toolCall entity.ToolCall) (*entity.ToolResult, error) {
	uc.mu.RLock()
	if uc.connection.Status != entity.ConnectionStatusConnected {
		uc.mu.RUnlock()
		return nil, fmt.Errorf("not connected to server")
	}
	uc.mu.RUnlock()

	result, err := uc.mcpRepo.CallTool(ctx, toolCall)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tool %s: %w", toolCall.Name, err)
	}

	log.Printf("Successfully executed tool: %s", toolCall.Name)
	return result, nil
}

// HandleIncomingMessage handles incoming messages from the server
func (uc *mcpUseCase) HandleIncomingMessage(ctx context.Context, message *entity.Message) error {
	uc.mu.RLock()
	handler, exists := uc.handlers[message.Method]
	uc.mu.RUnlock()

	if exists {
		return handler(message)
	}

	// Default handling for common methods
	switch message.Method {
	case "ping":
		return uc.handlePing(ctx, message)
	case "tools/list":
		return uc.handleToolsList(ctx, message)
	case "tools/call":
		return uc.handleToolsCall(ctx, message)
	default:
		log.Printf("Unhandled message method: %s", message.Method)
	}

	return nil
}

// SendOutgoingMessage sends a message to the server
func (uc *mcpUseCase) SendOutgoingMessage(ctx context.Context, message *entity.Message) error {
	uc.mu.RLock()
	if uc.connection.Status != entity.ConnectionStatusConnected {
		uc.mu.RUnlock()
		return fmt.Errorf("not connected to server")
	}
	uc.mu.RUnlock()

	return uc.mcpRepo.SendMessage(ctx, message)
}

// RegisterHandler registers a message handler for a specific method
func (uc *mcpUseCase) RegisterHandler(method string, handler MessageHandler) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.handlers[method] = handler
}

// handlePing handles ping messages
func (uc *mcpUseCase) handlePing(ctx context.Context, message *entity.Message) error {
	log.Println("Received ping, sending pong")
	pongMsg := &entity.Message{
		ID:     message.ID,
		Method: "pong",
	}
	return uc.SendOutgoingMessage(ctx, pongMsg)
}

// handleToolsList handles tools/list messages
func (uc *mcpUseCase) handleToolsList(ctx context.Context, message *entity.Message) error {
	log.Println("Received tools/list request")
	// Implementation would depend on specific requirements
	return nil
}

// handleToolsCall handles tools/call messages
func (uc *mcpUseCase) handleToolsCall(ctx context.Context, message *entity.Message) error {
	log.Println("Received tools/call request")
	// Implementation would depend on specific requirements
	return nil
}
