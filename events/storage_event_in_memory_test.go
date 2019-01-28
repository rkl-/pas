package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testEvent01 struct {
}

func (e *testEvent01) GetName() string {
	return "event.test_01"
}

type testEvent02 struct {
}

func (e *testEvent02) GetName() string {
	return "event.test_02"
}

// TestInMemoryEventStorage_AddEvent
//
//
func TestInMemoryEventStorage_AddEvent(t *testing.T) {
	getEventStorage(t)
}

// TestInMemoryEventStorage_GetEventStream
//
//
func TestInMemoryEventStorage_GetEventStream(t *testing.T) {
	storage := getEventStorage(t)

	fetchedStreams := []Event{}

	for event := range storage.GetEventStream() {
		fetchedStreams = append(fetchedStreams, event)
	}

	assert.Equal(t, storage.(*InMemoryEventStorage).events, fetchedStreams)
}

func getEventStorage(t *testing.T) EventStorage {
	storage := &InMemoryEventStorage{}

	// add event #1
	event01 := &testEvent01{}
	storage.AddEvent(event01)
	assert.Len(t, storage.events, 1)

	// add event #2
	event02 := &testEvent02{}
	storage.AddEvent(event02)
	assert.Len(t, storage.events, 2)

	return storage
}
