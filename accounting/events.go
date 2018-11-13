package accounting

import "github.com/satori/go.uuid"

var eventDispatcherInstance *EventDispatcher

// EventInterface event interface
//
//
type EventInterface interface {
	GetName() string
}

// EventSubscriberInterface interface for event subscriber
//
//
type EventHandlerInterface interface {
	Handle(event EventInterface)
}

// EventDispatcher accounting event dispatcher
//
//
type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

// GetInstance creates new event dispatcher
//
//
func (d EventDispatcher) GetInstance() *EventDispatcher {
	if eventDispatcherInstance == nil {
		ed := &EventDispatcher{
			handlers: map[string][]EventHandlerInterface{},
		}

		eventDispatcherInstance = ed
	}

	return eventDispatcherInstance
}

// RegisterHandler register an event handler
//
//
func (d *EventDispatcher) RegisterHandler(eventName string, handler EventHandlerInterface) {
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

// Dispatch dispatch an event
//
//
func (d *EventDispatcher) Dispatch(event EventInterface) {
	if handlers, ok := d.handlers[event.GetName()]; ok {
		for _, handler := range handlers {
			handler.Handle(event)
		}
	}
}

// AccountCreatedEvent event when an account was created
//
//
type AccountCreatedEvent struct {
	accountId uuid.UUID
}

// GetName get event name
//
//
func (e *AccountCreatedEvent) GetName() string {
	return "event.account_created"
}
