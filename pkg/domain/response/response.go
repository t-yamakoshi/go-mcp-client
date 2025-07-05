package response

import (
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
)

// InitializeResponse represents the initialize response
type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      entity.ServerInfo  `json:"serverInfo"`
}

// ServerCapabilities represents server capabilities
type ServerCapabilities struct {
	Tools map[string]interface{} `json:"tools,omitempty"`
}
