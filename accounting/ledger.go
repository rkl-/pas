package accounting

import "github.com/satori/go.uuid"

// Ledger ledger for accounting
//
//
type Ledger struct {
	eventDispatcher *EventDispatcher
}

// GetInstance create new ledger instance
//
//
func (l Ledger) New(eventDispatcher *EventDispatcher) *Ledger {
	le := &Ledger{
		eventDispatcher: eventDispatcher,
	}

	return le
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

	l.eventDispatcher.Dispatch(&AccountCreatedEvent{a.id})

	return a
}
