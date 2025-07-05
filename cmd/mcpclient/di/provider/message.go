package provider

import (
	"github.com/google/wire"
	"github.com/t-yamakoshi/go-mcp-client/pkg/interfaces/message"
)

var MessageSet = wire.NewSet(
	message.NewMessageHandler,
)
