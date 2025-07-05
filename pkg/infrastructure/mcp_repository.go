package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/response"
)

var _ IFMCPRepository = (*MCPRepositoryImpl)(nil)

type IFMCPRepository interface {
	Connect(ctx context.Context, serverURL string) error
	Disconnect() error
	IsConnected() bool
	SendMessage(ctx context.Context, message *entity.Message) error
	ReceiveMessage(ctx context.Context) (*entity.Message, error)
}

// MCPRepositoryImpl implements the MCP repository interface
type MCPRepositoryImpl struct {
	conn     *websocket.Conn
	mu       sync.RWMutex
	handlers map[string]MessageHandler
}

// MessageHandler is a function type for handling incoming messages
type MessageHandler func(*entity.Message) error

// NewMCPRepositoryImpl creates a new MCP repository implementation
func NewMCPRepositoryImpl() *MCPRepositoryImpl {
	return &MCPRepositoryImpl{
		handlers: make(map[string]MessageHandler),
	}
}

// Connect establishes a WebSocket connection to the MCP server
func (r *MCPRepositoryImpl) Connect(ctx context.Context, serverURL string) error {
	u, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	r.conn = conn

	// Start listening for messages
	go r.listen()

	return nil
}

// Disconnect closes the WebSocket connection
func (r *MCPRepositoryImpl) Disconnect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// IsConnected returns whether the client is connected
func (r *MCPRepositoryImpl) IsConnected() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.conn != nil
}

// SendMessage sends a message to the server
func (r *MCPRepositoryImpl) SendMessage(ctx context.Context, message *entity.Message) error {
	r.mu.RLock()
	if r.conn == nil {
		r.mu.RUnlock()
		return fmt.Errorf("not connected")
	}
	r.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.conn.WriteMessage(websocket.TextMessage, data)
}

// ReceiveMessage receives a message from the server
func (r *MCPRepositoryImpl) ReceiveMessage(ctx context.Context) (*entity.Message, error) {
	r.mu.RLock()
	if r.conn == nil {
		r.mu.RUnlock()
		return nil, fmt.Errorf("not connected")
	}
	r.mu.RUnlock()

	// Set read deadline
	if err := r.conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	_, data, err := r.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	var msg entity.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &msg, nil
}

// Initialize sends an initialize request to the server
func (r *MCPRepositoryImpl) Initialize(ctx context.Context, clientInfo entity.ClientInfo) (*response.InitializeResponse, error) {
	req := struct {
		ProtocolVersion string                 `json:"protocolVersion"`
		Capabilities    map[string]interface{} `json:"capabilities"`
		ClientInfo      entity.ClientInfo      `json:"clientInfo"`
	}{
		ProtocolVersion: "2024-11-05",
		Capabilities:    make(map[string]interface{}),
		ClientInfo:      clientInfo,
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal initialize request: %w", err)
	}

	msg := &entity.Message{
		ID:     uuid.New().String(),
		Method: "initialize",
		Params: reqData,
	}

	// Send the message
	if err := r.SendMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to send initialize message: %w", err)
	}

	// Wait for response (simplified implementation)
	time.Sleep(100 * time.Millisecond)

	// Return a mock response for now
	return &response.InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: response.ServerCapabilities{
			Tools: make(map[string]interface{}),
		},
		ServerInfo: entity.ServerInfo{
			Name:    "mock-mcp-server",
			Version: "1.0.0",
		},
	}, nil
}

// ListTools retrieves available tools from the server
func (r *MCPRepositoryImpl) ListTools(ctx context.Context) ([]entity.Tool, error) {
	msg := &entity.Message{
		ID:     uuid.New().String(),
		Method: "tools/list",
	}

	// Send the message
	if err := r.SendMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to send tools/list message: %w", err)
	}

	// Return empty list for now
	return []entity.Tool{}, nil
}

// CallTool executes a tool on the server
func (r *MCPRepositoryImpl) CallTool(ctx context.Context, toolCall entity.ToolCall) (*entity.ToolResult, error) {
	reqData, err := json.Marshal(toolCall)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tool call: %w", err)
	}

	msg := &entity.Message{
		ID:     uuid.New().String(),
		Method: "tools/call",
		Params: reqData,
	}

	// Send the message
	if err := r.SendMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to send tools/call message: %w", err)
	}

	// Return empty result for now
	return &entity.ToolResult{
		Content: []entity.Content{},
		IsError: false,
	}, nil
}

// RegisterHandler registers a message handler for a specific method
func (r *MCPRepositoryImpl) RegisterHandler(method string, handler MessageHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[method] = handler
}

// listen listens for incoming messages
func (r *MCPRepositoryImpl) listen() {
	for {
		r.mu.RLock()
		if r.conn == nil {
			r.mu.RUnlock()
			return
		}
		r.mu.RUnlock()

		_, data, err := r.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		var msg entity.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Handle the message
		if err := r.handleMessage(&msg); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}

// handleMessage handles an incoming message
func (r *MCPRepositoryImpl) handleMessage(msg *entity.Message) error {
	r.mu.RLock()
	handler, exists := r.handlers[msg.Method]
	r.mu.RUnlock()

	if exists {
		return handler(msg)
	}

	// Default handling for common methods
	switch msg.Method {
	case "initialize":
		return r.handleInitialize(msg)
	case "tools/list":
		return r.handleToolsList(msg)
	case "tools/call":
		return r.handleToolsCall(msg)
	default:
		log.Printf("Unhandled method: %s", msg.Method)
	}

	return nil
}

// handleInitialize handles initialize responses
func (r *MCPRepositoryImpl) handleInitialize(msg *entity.Message) error {
	log.Printf("Received initialize response")
	return nil
}

// handleToolsList handles tools/list responses
func (r *MCPRepositoryImpl) handleToolsList(msg *entity.Message) error {
	log.Printf("Received tools list")
	return nil
}

// handleToolsCall handles tools/call requests
func (r *MCPRepositoryImpl) handleToolsCall(msg *entity.Message) error {
	log.Printf("Received tool call request")
	return nil
}
