package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
	"pas/events"
)

// Ledger ledger for accounting
//
//
type Ledger struct {
	eventDispatcher   events.EventDispatcher
	accountRepository *AccountRepository
}

// GetInstance create new ledger instance
//
//
func (l Ledger) New(eventDispatcher events.EventDispatcher, eventStorage events.EventStorage) *Ledger {
	if eventDispatcher == nil {
		panic("event dispatcher is required")
	}

	le := &Ledger{
		eventDispatcher:   eventDispatcher,
		accountRepository: &AccountRepository{eventStorage},
	}

	return le
}

// CreateAccount create a new ledger account and dispatch AccountCreatedEvent
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
		currencyId:   currencyId,
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

// AddValue add new value to an account and dispatch AccountValueAddedEvent
//
//
func (l *Ledger) AddValue(account *Account, value Money, reason string) error {
	if err := account.addValue(value, reason); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(&AccountValueAddedEvent{
		accountId: account.id,
		value:     value,
		reason:    reason,
	})

	return nil
}

// SubtractValue subtract value from an account and dispatch AccountValueSubtractedEvent
//
//
func (l *Ledger) SubtractValue(account *Account, value Money, reason string) error {
	if err := account.subtractValue(value, reason); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(&AccountValueSubtractedEvent{
		accountId: account.id,
		value:     value,
		reason:    reason,
	})

	return nil
}

// LoadAccount an account by id
//
//
func (l *Ledger) LoadAccount(accountId uuid.UUID) (*Account, error) {
	return l.accountRepository.loadById(accountId)
}
