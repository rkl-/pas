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

// RequestHandler request handler
//
//
type RequestHandler interface {
	Handle(request interface{}) (interface{}, error)
}

// RequestBus generic request bus
//
//
type RequestBus interface {
	RegisterHandler(requestId string, handler RequestHandler) error
	Execute(request interface{}) (interface{}, error)
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

	return nil
}
