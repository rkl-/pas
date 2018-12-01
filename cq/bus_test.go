package cq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testRequestHandler struct {
}

func (h *testRequestHandler) Handle(request interface{}) (interface{}, error) {
	return nil, nil
}

// TestCommandBus_RegisterHandler
//
//
func TestCommandBus_RegisterHandler(t *testing.T) {
	// get command bus instance
	bus := CommandBus{}.GetInstance()

	// negative test for InvalidHandlerIdError
	bus.handlers = map[string]RequestHandler{
		"test-handler-id": nil,
	}

	err := bus.RegisterHandler("test-handler-id", nil)
	assert.IsType(t, &InvalidHandlerIdError{}, err)

	// negative test for HandlerAlreadyRegisteredError
	bus.handlers = map[string]RequestHandler{
		"command.test-handler-id": nil,
	}

	err = bus.RegisterHandler("command.test-handler-id", nil)
	assert.IsType(t, &HandlerAlreadyRegisteredError{}, err)

	// positive test
	err = bus.RegisterHandler("command.test_command", &testRequestHandler{})
	assert.Nil(t, err)
}
