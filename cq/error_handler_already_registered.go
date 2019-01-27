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
