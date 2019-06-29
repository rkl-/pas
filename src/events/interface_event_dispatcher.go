package events

// EventDispatcher event dispatcher
//
//
type EventDispatcher interface {
	RegisterHandler(eventName string, handler EventHandler)
	Dispatch(event Event)
}
