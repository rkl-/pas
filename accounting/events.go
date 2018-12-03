package accounting

import (
	"github.com/satori/go.uuid"
	"pas/events"
)

// EventStorage common storage for events
//
//
type EventStorage interface {
	AddEvent(event events.Event)
	GetEventStream() chan events.Event
}

// SingleAccountEvent event which has an unique account association
//
//
type SingleAccountEvent interface {
	GetAccountId() uuid.UUID
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
