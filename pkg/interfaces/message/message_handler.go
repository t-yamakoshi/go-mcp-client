package message

import (
	"context"
	"log"

	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
	"github.com/t-yamakoshi/go-mcp-client/pkg/usecase"
)

// MessageHandler defines the interface for MCP message handling operations

var _ IFMessageHandler = (*MessageHandler)(nil)

type IFMessageHandler interface {
	RegisterHandlers(mcpUsecase *usecase.IFMCPUsecase)
}

// MessageHandler implements the MessageHandler interface
type MessageHandler struct{}

// NewMessageHandler creates a new message handler
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

// RegisterHandlers registers message handlers for the MCP client
func (h *MessageHandler) RegisterHandlers(mcpUsecase *usecase.IFMCPUsecase) {
	// Register handler for tools/list
	(*mcpUsecase).RegisterHandler("tools/list", func(msg *entity.Message) error {
		log.Println("Received tools/list request")
		// Here you would implement the actual tools listing logic
		return nil
	})

	// Register handler for tools/call
	(*mcpUsecase).RegisterHandler("tools/call", func(msg *entity.Message) error {
		log.Println("Received tools/call request")
		// Here you would implement the actual tool calling logic
		return nil
	})

	// Register handler for ping
	(*mcpUsecase).RegisterHandler("ping", func(msg *entity.Message) error {
		log.Println("Received ping, sending pong")
		// Send pong response
		pongMsg := &entity.Message{
			ID:     msg.ID,
			Method: "pong",
		}
		return (*mcpUsecase).SendOutgoingMessage(context.Background(), pongMsg)
	})
}
