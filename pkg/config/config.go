package config

import (
	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
)

// Config represents application configuration
type Config struct {
	ServerURL  string            `json:"server_url"`
	ClientInfo entity.ClientInfo `json:"client_info"`
	LogLevel   string            `json:"log_level"`
}
