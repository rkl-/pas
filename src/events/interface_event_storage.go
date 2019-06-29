package events

// EventStorage common storage for events
//
//
type EventStorage interface {
	AddEvent(event Event)
	GetEventStream() chan Event
}
