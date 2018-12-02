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

// RequestHandlerNotRegisteredError
//
//
type RequestHandlerNotRegisteredError struct {
	requestId string
}

func (e *RequestHandlerNotRegisteredError) Error() string {
	return fmt.Sprintf("no handler registered for request id \"%s\"", e.requestId)
}

// RequestNotSupportedError
//
//
type UnsupportedRequestError struct {
}

func (e *UnsupportedRequestError) Error() string {
	return fmt.Sprintf("unsupported request")
}
