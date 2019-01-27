package accounting

import "github.com/satori/go.uuid"

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
