package entity

import (
	"time"
)

// Connection represents a connection state
type Connection struct {
	ID        string
	ServerURL string
	Status    ConnectionStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ConnectionStatus represents the status of a connection
type ConnectionStatus string

const (
	ConnectionStatusDisconnected ConnectionStatus = "disconnected"
	ConnectionStatusConnecting   ConnectionStatus = "connecting"
	ConnectionStatusConnected    ConnectionStatus = "connected"
	ConnectionStatusError        ConnectionStatus = "error"
)
