package message

import "github.com/google/wire"

var MessageSet = wire.NewSet(
	NewMessageHandler,
)
