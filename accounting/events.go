package accounting

import "github.com/satori/go.uuid"

var eventDispatcherInstance *EventDispatcher

// EventStorage common storage for events
//
//
type EventStorage interface {
	AddEvent(event Event)
	GetEventStream() chan Event
}

// Event event interface
//
//
type Event interface {
	GetName() string
}

// SingleAccountEvent event which has an unique account association
//
//
type SingleAccountEvent interface {
	GetAccountId() uuid.UUID
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

////////////////////////////////////////////////////////////////////
// EVENTS BELOW
////////////////////////////////////////////////////////////////////

// AccountCreatedEvent event when an accountId was created
//
//
type AccountCreatedEvent struct {
	accountId    uuid.UUID
	accountTitle string
	currencyId   string
}

func (e *AccountCreatedEvent) GetName() string {
	return "event.account_created"
}

func (e *AccountCreatedEvent) GetAccountId() uuid.UUID {
	return e.accountId
}

// AccountValueTransferredEvent event when value was transferred fromId one accountId toId another
//
//
type AccountValueTransferredEvent struct {
	fromId uuid.UUID
	toId   uuid.UUID
	value  Money
	reason string
}

func (e *AccountValueTransferredEvent) GetName() string {
	return "event.account_value_transferred"
}

// AccountValueAddedEvent event when new value was added to an accountId
//
//
type AccountValueAddedEvent struct {
	accountId uuid.UUID
	value     Money
	reason    string
}

func (e *AccountValueAddedEvent) GetName() string {
	return "event.account_value_added"
}

func (e *AccountValueAddedEvent) GetAccountId() uuid.UUID {
	return e.accountId
}

// AccountValueSubtractedEvent event when value from an account was subtracted
//
//
type AccountValueSubtractedEvent struct {
	accountId uuid.UUID
	value     Money
	reason    string
}

func (e *AccountValueSubtractedEvent) GetName() string {
	return "event.account_value_subtracted"
}

func (e *AccountValueSubtractedEvent) GetAccountId() uuid.UUID {
	return e.accountId
}
