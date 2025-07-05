package cli

import "github.com/google/wire"

var CLISet = wire.NewSet(
	NewCLIHandler,
)
