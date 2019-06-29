package events

import (
	"github.com/satori/go.uuid"
	"pas/src/money"
)

// AccountValueSubtractedEvent event when Value from an account was subtracted
//
//
type AccountValueSubtractedEvent struct {
	AccountId uuid.UUID
	Value     money.Money
	Reason    string
}

func (e *AccountValueSubtractedEvent) GetName() string {
	return "event.account_value_subtracted"
}

func (e *AccountValueSubtractedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
