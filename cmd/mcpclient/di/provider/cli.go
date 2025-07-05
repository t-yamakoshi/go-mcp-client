package provider

import (
	"github.com/google/wire"
	"github.com/t-yamakoshi/go-mcp-client/pkg/interfaces/cli"
)

var CLISet = wire.NewSet(
	cli.NewCLIHandler,
)
