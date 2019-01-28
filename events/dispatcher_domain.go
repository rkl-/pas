package events

// DomainDispatcher accounting event dispatcher
//
//
type DomainDispatcher struct {
	handlers map[string][]EventHandler
}

// New creates new event dispatcher
//
//
func (d DomainDispatcher) New() EventDispatcher {
	dd := &DomainDispatcher{
		handlers: map[string][]EventHandler{},
	}

	return dd
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
