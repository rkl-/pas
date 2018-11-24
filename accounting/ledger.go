package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
	"strings"
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

// Account a ledger accountId
//
//
type Account struct {
	id      uuid.UUID
	title   string
	balance Money
}

// CreateAccount create a new ledger accountId
//
//
func (l *Ledger) CreateAccount(title, currencyId string) *Account {
	a := &Account{
		id:      uuid.NewV4(),
		title:   title,
		balance: Money{}.NewFromInt(0, currencyId),
	}

	l.eventDispatcher.Dispatch(&AccountCreatedEvent{
		accountId:    a.id,
		accountTitle: title,
	})

	return a
}

// TransferValue transfer value fromId one accountId toId another
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
		fromId: fromAccount.id,
		toId:   toAccount.id,
		value:  value,
		reason: reason,
	})

	return nil
}

// AddValue add new value to an account
//
//
func (l *Ledger) AddValue(toAccount *Account, value Money, reason string) error {
	if strings.Compare(toAccount.balance.currencyId, value.currencyId) != 0 {
		return &UnequalCurrenciesError{}
	}

	currencyId := toAccount.balance.currencyId
	newAmount := (&big.Int{}).Add(toAccount.balance.amount, value.amount)

	newBalance, err := Money{}.NewFromString(newAmount.String(), currencyId)
	if err != nil {
		return err
	}

	toAccount.balance = newBalance

	l.eventDispatcher.Dispatch(&AccountValueAddedEvent{
		accountId: toAccount.id,
		value:     value,
		reason:    reason,
	})

	return nil
}

// SubtractValue subtract value from an account
//
//
func (l *Ledger) SubtractValue(fromAccount *Account, value Money, reason string) error {
	ok, err := fromAccount.balance.IsLowerThan(value)
	if err != nil {
		return err
	}
	if ok {
		return &InsufficientFoundsError{}
	}

	currencyId := fromAccount.balance.currencyId
	newAmount := (&big.Int{}).Sub(fromAccount.balance.amount, value.amount)

	newBalance, err := Money{}.NewFromString(newAmount.String(), currencyId)
	if err != nil {
		return err
	}

	fromAccount.balance = newBalance

	l.eventDispatcher.Dispatch(&AccountValueSubtractedEvent{
		accountId: fromAccount.id,
		value:     value,
		reason:    reason,
	})

	return nil
}
