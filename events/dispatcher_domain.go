package events

var eventDispatcherInstance *DomainDispatcher

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
