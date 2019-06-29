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

// TestEventDispatcher_RegisterHandler
//
//
func TestEventDispatcher_RegisterHandler(t *testing.T) {
	dispatcher := DomainDispatcher{}.New()
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
	dispatcher := DomainDispatcher{}.New()
	dispatcher.RegisterHandler("event.test", &testEventHandler{})

	testEventHandlerExecuted = false

	dispatcher.Dispatch(&testEvent{})

	assert.True(t, testEventHandlerExecuted)
}
