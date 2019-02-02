package accounting

import "github.com/satori/go.uuid"

// AccountValueAddedEvent event when new value was added to an account
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
