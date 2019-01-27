package cq

import "fmt"

// InvalidHandlerIdError
//
//
type InvalidHandlerIdError struct {
	requiredPrefix string
}

func (e *InvalidHandlerIdError) Error() string {
	return fmt.Sprintf("the handler id must have the prefix \"%s\"", e.requiredPrefix)
}
