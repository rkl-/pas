package cq

import "fmt"

// HandlerAlreadyRegisteredError
//
//
type HandlerAlreadyRegisteredError struct {
	handlerId string
}

func (e *HandlerAlreadyRegisteredError) Error() string {
	return fmt.Sprintf("a handler for \"%s\" is already registered", e.handlerId)
}

// InvalidHandlerIdError
//
//
type InvalidHandlerIdError struct {
	requiredPrefix string
}

func (e *InvalidHandlerIdError) Error() string {
	return fmt.Sprintf("the handler id must have the prefix \"%s\"", e.requiredPrefix)
}
