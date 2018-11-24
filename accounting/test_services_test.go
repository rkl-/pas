package accounting

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

	assert.Equal(t, storage.(*inMemoryEventStorage).events, fetchedStreams)
}

func getEventStorage(t *testing.T) EventStorage {
	storage := &inMemoryEventStorage{}

	// add event #1
	event01 := &AccountCreatedEvent{accountId: uuid.NewV4(), accountTitle: "test account"}
	storage.AddEvent(event01)
	assert.Len(t, storage.events, 1)

	// add event #2
	event02 := &AccountValueAddedEvent{
		accountId: uuid.NewV4(),
		value:     Money{}.NewFromInt(10000, "EUR"),
		reason:    "no reason",
	}
	storage.AddEvent(event02)
	assert.Len(t, storage.events, 2)

	return storage
}
