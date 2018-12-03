package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testEventHandlerExecuted = false

type testEvent struct {
}

func (e *testEvent) GetName() string {
	return "event.test"
}

type testEventHandler struct {
}

func (h *testEventHandler) Handle(event Event) {
	testEventHandlerExecuted = true
}

// TestEventDispatcher_GetInstance
//
//
func TestEventDispatcher_GetInstance(t *testing.T) {
	dispatcher01 := DomainDispatcher{}.GetInstance()
	dispatcher01.RegisterHandler("event.test", &testEventHandler{})

	dispatcher02 := DomainDispatcher{}.GetInstance()

	assert.Equal(t, dispatcher01, dispatcher02)
}

// TestEventDispatcher_RegisterHandler
//
//
func TestEventDispatcher_RegisterHandler(t *testing.T) {
	eventDispatcherInstance = nil

	dispatcher := DomainDispatcher{}.GetInstance()
	dispatcher.RegisterHandler("event.test", &testEventHandler{})

	domainDispatcher, ok := dispatcher.(*DomainDispatcher)
	assert.True(t, ok)

	assert.Equal(t, 1, len(domainDispatcher.handlers["event.test"]))

	foundType := false

	switch domainDispatcher.handlers["event.test"][0].(type) {
	case *testEventHandler:
		foundType = true
	}

	assert.True(t, foundType)
}

// TestEventDispatcher_Dispatch
//
//
func TestEventDispatcher_Dispatch(t *testing.T) {
	eventDispatcherInstance = nil

	dispatcher := DomainDispatcher{}.GetInstance()
	dispatcher.RegisterHandler("event.test", &testEventHandler{})

	testEventHandlerExecuted = false

	dispatcher.Dispatch(&testEvent{})

	assert.True(t, testEventHandlerExecuted)
}
