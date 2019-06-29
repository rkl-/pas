package cq

import "fmt"

// RequestHandlerNotRegisteredError
//
//
type RequestHandlerNotRegisteredError struct {
	requestId string
}

func (e *RequestHandlerNotRegisteredError) Error() string {
	return fmt.Sprintf("no handler registered for request id \"%s\"", e.requestId)
}
