package accounting

import "github.com/satori/go.uuid"

// AccountValueTransferredEvent event when value was transferred fromId one AccountId toId another
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
