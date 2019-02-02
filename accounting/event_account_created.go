package accounting

import "github.com/satori/go.uuid"

// AccountCreatedEvent event when an account was created
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
