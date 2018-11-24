package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
)

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

// TransferValue transfer value from one account to another
//
//
func (l *Ledger) TransferValue(fromAccount, toAccount *Account, value Money, reason string) error {
	ok, err := fromAccount.balance.IsLowerThan(value)
	if err != nil {
		return err
	}
	if ok {
		return &InsufficientFoundsError{}
	}

	currencyId := fromAccount.balance.GetCurrencyId()
	oldFromAmount := fromAccount.balance.GetAmount()
	oldToAmount := toAccount.balance.GetAmount()

	newFromAmount := (&big.Int{}).Sub(oldFromAmount, value.GetAmount())
	newToAmount := (&big.Int{}).Add(oldToAmount, value.GetAmount())

	newFromBalance, err := Money{}.NewFromString(newFromAmount.String(), currencyId)
	if err != nil {
		return err
	}
	newToBalance, err := Money{}.NewFromString(newToAmount.String(), currencyId)
	if err != nil {
		return err
	}

	fromAccount.balance = newFromBalance
	toAccount.balance = newToBalance

	l.eventDispatcher.Dispatch(&AccountValueTransferredEvent{
		from: fromAccount.id,
		to:   toAccount.id,
	})

	return nil
}
