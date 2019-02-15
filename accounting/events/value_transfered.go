package events

import (
	"github.com/satori/go.uuid"
	"pas/money"
)

// AccountValueTransferredEvent event when Value was transferred from one account to another
//
//
type AccountValueTransferredEvent struct {
	FromId uuid.UUID
	ToId   uuid.UUID
	Value  money.Money
	Reason string
}

func (e *AccountValueTransferredEvent) GetName() string {
	return "event.account_value_transferred"
}
