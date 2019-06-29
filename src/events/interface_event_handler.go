package events

// EventHandler interface for event subscriber
//
//
type EventHandler interface {
	Handle(event Event)
}
