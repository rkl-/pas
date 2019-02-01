package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
	"pas/events"
)

// DefaultLedger ledger for accounting
//
//
type DefaultLedger struct {
	eventDispatcher   events.EventDispatcher
	accountRepository *AccountRepository
}

// New create new ledger instance
//
//
func (l DefaultLedger) New(eventDispatcher events.EventDispatcher, eventStorage events.EventStorage) Ledger {
	if eventDispatcher == nil {
		panic("event dispatcher is required")
	}

	le := &DefaultLedger{
		eventDispatcher:   eventDispatcher,
		accountRepository: &AccountRepository{eventStorage},
	}

	return le
}

// CreateAccount create a new ledger account and dispatch AccountCreatedEvent
//
//
func (l *DefaultLedger) CreateAccount(title, currencyId string) (*Account, error) {
	a := &Account{
		id:      uuid.NewV4(),
		title:   title,
		balance: Money{}.NewFromInt(0, currencyId),
	}

	event := &AccountCreatedEvent{
		accountId:    a.id,
		accountTitle: title,
		currencyId:   currencyId,
	}

	a.addRecordedEvent(event)

	if err := l.accountRepository.save(a); err != nil {
		return nil, err
	}

	l.eventDispatcher.Dispatch(event)

	return a, nil
}

// TransferValue transfer value fromId one AccountId toId another
//
//
func (l *DefaultLedger) TransferValue(fromAccount, toAccount *Account, value Money, reason string) error {
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

	event := &AccountValueTransferredEvent{
		fromId: fromAccount.id,
		toId:   toAccount.id,
		value:  value,
		reason: reason,
	}

	fromAccount.addRecordedEvent(event)

	if err := l.accountRepository.save(fromAccount); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(event)

	return nil
}

// AddValue add new value to an account and dispatch AccountValueAddedEvent
//
//
func (l *DefaultLedger) AddValue(account *Account, value Money, reason string) error {
	if err := account.addValue(value, reason); err != nil {
		return err
	}

	event := &AccountValueAddedEvent{
		accountId: account.id,
		value:     value,
		reason:    reason,
	}

	account.addRecordedEvent(event)

	if err := l.accountRepository.save(account); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(event)

	return nil
}

// SubtractValue subtract value from an account and dispatch AccountValueSubtractedEvent
//
//
func (l *DefaultLedger) SubtractValue(account *Account, value Money, reason string) error {
	if err := account.subtractValue(value, reason); err != nil {
		return err
	}

	event := &AccountValueSubtractedEvent{
		accountId: account.id,
		value:     value,
		reason:    reason,
	}

	account.addRecordedEvent(event)

	if err := l.accountRepository.save(account); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(event)

	return nil
}

// AddPlannedCashReceipt add a planned cash receipt to an account
//
//
func (l *DefaultLedger) AddPlannedCashReceipt(account *Account, receipt *PlannedCashReceipt) error {
	if err := account.addPlannedCashReceipt(receipt); err != nil {
		return err
	}

	event := PlannedCashReceiptCreatedEvent{}.New(
		account.GetId(),
		receipt.date,
		receipt.amount,
		receipt.title,
	)

	account.addRecordedEvent(event)

	if err := l.accountRepository.save(account); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(event)

	return nil
}

// HasAccount efficient way to check if an account exists or not.
//
//
func (l *DefaultLedger) HasAccount(accountId uuid.UUID) bool {
	return l.accountRepository.hasAccount(accountId)
}

// LoadAccount an account by id
//
//
func (l *DefaultLedger) LoadAccount(accountId uuid.UUID) (*Account, error) {
	return l.accountRepository.loadById(accountId)
}
