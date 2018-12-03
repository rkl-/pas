package events

var eventDispatcherInstance *EventDispatcher

// Event event interface
//
//
type Event interface {
	GetName() string
}

// EventHandler interface for event subscriber
//
//
type EventHandler interface {
	Handle(event Event)
}

// EventDispatcher accounting event dispatcher
//
//
type EventDispatcher struct {
	handlers map[string][]EventHandler
}

// GetInstance creates new event dispatcher
//
//
func (d EventDispatcher) GetInstance() *EventDispatcher {
	if eventDispatcherInstance == nil {
		ed := &EventDispatcher{
			handlers: map[string][]EventHandler{},
		}

		eventDispatcherInstance = ed
	}

	return eventDispatcherInstance
}

// RegisterHandler register an event handler
//
//
func (d *EventDispatcher) RegisterHandler(eventName string, handler EventHandler) {
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

// Dispatch dispatch an event
//
//
func (d *EventDispatcher) Dispatch(event Event) {
	if handlers, ok := d.handlers[event.GetName()]; ok {
		for _, handler := range handlers {
			handler.Handle(event)
		}
	}
}
