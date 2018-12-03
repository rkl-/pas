package events

var eventDispatcherInstance *DomainDispatcher

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

// EventDispatcher event dispatcher
//
//
type EventDispatcher interface {
	RegisterHandler(eventName string, handler EventHandler)
	Dispatch(event Event)
}

// DomainDispatcher accounting event dispatcher
//
//
type DomainDispatcher struct {
	handlers map[string][]EventHandler
}

// GetInstance creates new event dispatcher
//
//
func (d DomainDispatcher) GetInstance() EventDispatcher {
	if eventDispatcherInstance == nil {
		ed := &DomainDispatcher{
			handlers: map[string][]EventHandler{},
		}

		eventDispatcherInstance = ed
	}

	return eventDispatcherInstance
}

// RegisterHandler register an event handler
//
//
func (d *DomainDispatcher) RegisterHandler(eventName string, handler EventHandler) {
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

// Dispatch dispatch an event
//
//
func (d *DomainDispatcher) Dispatch(event Event) {
	if handlers, ok := d.handlers[event.GetName()]; ok {
		for _, handler := range handlers {
			handler.Handle(event)
		}
	}
}
