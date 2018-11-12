package accounting

import "github.com/satori/go.uuid"

// Ledger ledger for accounting
//
//
type Ledger struct {
}

// Account a ledger account
//
//
type Account struct {
	id      uuid.UUID
	title   string
	balance Money
}

// CreateAccount create a new ledger account
//
//
func (l *Ledger) CreateAccount(title, currencyId string) *Account {
	a := &Account{
		id:      uuid.NewV4(),
		title:   title,
		balance: Money{}.NewFromInt(0, currencyId),
	}

	// TODO, dispatch domain event: AccountCreatedEvent

	return a
}
