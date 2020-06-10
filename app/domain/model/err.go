package model

import (
	"fmt"
)

var (
	ErrNotFound      = fmt.Errorf("not found")
	ErrInvalidAmount = fmt.Errorf("invalid amount")
	ErrInvalid       = fmt.Errorf("invalid")
)
