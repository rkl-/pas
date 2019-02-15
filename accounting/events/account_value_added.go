package events

import (
	"github.com/satori/go.uuid"
	"pas/money"
)

// AccountValueAddedEvent event when new Value was added to an account
//
//
type AccountValueAddedEvent struct {
	AccountId uuid.UUID
	Value     money.Money
	Reason    string
}

func (e *AccountValueAddedEvent) GetName() string {
	return "event.account_value_added"
}

func (e *AccountValueAddedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
