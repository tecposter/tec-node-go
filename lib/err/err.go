package err

import (
	"errors"
)

// Global errors
var (
	ErrCmdNotFound = errors.New("Command not found in post module")
)
