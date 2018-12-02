package cq

import "strings"

var commandBusInstance *CommandBus

// CommandBus command bus
//
//
type CommandBus struct {
	genericRequestBus
}

func (b CommandBus) GetInstance() *CommandBus {
	if commandBusInstance == nil {
		commandBusInstance = &CommandBus{}
		commandBusInstance.genericRequestBus.handlerPrefix = "command."
	}

	return commandBusInstance
}

// Request any kind of request
//
//
type Request interface {
	GetRequestId() string
}

// RequestHandler request handler
//
//
type RequestHandler interface {
	Handle(request Request) (interface{}, error)
}

// RequestBus generic request bus
//
//
type RequestBus interface {
	RegisterHandler(requestId string, handler RequestHandler) error
	Execute(request Request) (interface{}, error)
}

// genericRequestBus generic request bus
//
//
type genericRequestBus struct {
	handlerPrefix string
	handlers      map[string]RequestHandler
}

func (b *genericRequestBus) RegisterHandler(requestId string, handler RequestHandler) error {
	if b.handlers == nil {
		b.handlers = map[string]RequestHandler{}
	}

	if !strings.HasPrefix(requestId, b.handlerPrefix) {
		return &InvalidHandlerIdError{b.handlerPrefix}
	}

	if _, ok := b.handlers[requestId]; ok {
		return &HandlerAlreadyRegisteredError{requestId}
	}

	b.handlers[requestId] = handler

	return nil
}

func (b *genericRequestBus) Execute(request Request) (interface{}, error) {
	requestId := request.GetRequestId()

	if _, ok := b.handlers[requestId]; !ok {
		return nil, &RequestHandlerNotRegisteredError{requestId}
	}

	return b.handlers[requestId].Handle(request)
}
