package cq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type unregisteredTestCommand struct {
}

func (c *unregisteredTestCommand) GetRequestId() string {
	return "command.unregistered_test_command"
}

type registeredTestCommand struct {
	value string
}

func (c *registeredTestCommand) GetRequestId() string {
	return "command.registered_test_command"
}

type testRequestHandler struct {
}

func (h *testRequestHandler) Handle(request Request) (interface{}, error) {
	command, _ := request.(*registeredTestCommand)

	return command.value, nil
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
	assert.Len(t, bus.handlers, 2)
}

// TestGenericRequestBus_Execute
//
//
func TestGenericRequestBus_Execute(t *testing.T) {
	// get command bus instance
	bus := CommandBus{}.GetInstance()

	err := bus.RegisterHandler("command.registered_test_command", &testRequestHandler{})
	assert.Nil(t, err)

	// negative test for RequestHandlerNotRegisteredError
	_, err = bus.Execute(&unregisteredTestCommand{})
	assert.IsType(t, &RequestHandlerNotRegisteredError{}, err)

	// positive test
	command := &registeredTestCommand{"test-value"}

	res, err := bus.Execute(command)
	assert.Nil(t, err)

	strRes, ok := res.(string)
	assert.True(t, ok)
	assert.Equal(t, "test-value", strRes)
}
