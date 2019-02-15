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

	return a, l.addAndDispatchAccountEvent(a, event)
}

// TransferValue transfer value from one Account to another
//
//
func (l *DefaultLedger) TransferValue(fromAccountId, toAccountId uuid.UUID, value Money, reason string) error {
	fromAccount, err := l.LoadAccount(fromAccountId)
	if err != nil {
		return err
	}

	toAccount, err := l.LoadAccount(toAccountId)
	if err != nil {
		return err
	}

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

	return l.addAndDispatchAccountEvent(fromAccount, event)
}

// AddValue add new value to an account and dispatch AccountValueAddedEvent
//
//
func (l *DefaultLedger) AddValue(accountId uuid.UUID, value Money, reason string) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	if err := account.addValue(value, reason); err != nil {
		return err
	}

	event := &AccountValueAddedEvent{
		accountId: accountId,
		value:     value,
		reason:    reason,
	}

	return l.addAndDispatchAccountEvent(account, event)
}

// SubtractValue subtract value from an account and dispatch AccountValueSubtractedEvent
//
//
func (l *DefaultLedger) SubtractValue(accountId uuid.UUID, value Money, reason string) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	if err := account.subtractValue(value, reason); err != nil {
		return err
	}

	event := &AccountValueSubtractedEvent{
		accountId: accountId,
		value:     value,
		reason:    reason,
	}

	return l.addAndDispatchAccountEvent(account, event)
}

// AddPlannedCashReceipt add a planned cash receipt to an account
//
//
func (l *DefaultLedger) AddPlannedCashReceipt(accountId uuid.UUID, receipt *PlannedCashFlow) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	if err := account.addPlannedCashReceipt(receipt); err != nil {
		return err
	}

	event := PlannedCashReceiptCreatedEvent{}.NewFrom(receipt)

	return l.addAndDispatchAccountEvent(account, event)
}

// ConfirmPlannedCashReceipt confirm a planned cash receipt
//
//
func (l *DefaultLedger) ConfirmPlannedCashReceipt(accountId uuid.UUID, receiptId uuid.UUID) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	flows := account.GetPlannedCashReceipts()
	receipt, ok := flows[receiptId]
	if !ok {
		return &PlannedCashReceiptNotFoundError{receiptId, accountId}
	}

	if err := account.addValue(receipt.amount, receipt.title); err != nil {
		return err
	}

	delete(account.plannedCashReceipts, receiptId)

	event := PlannedCashReceiptConfirmedEvent{}.NewFrom(receipt)

	if err := l.addAndDispatchAccountEvent(account, event); err != nil {
		return err
	}

	return nil
}

// AddPlannedCashWithdrawal add a planned cash withdrawal to an account
//
//
func (l *DefaultLedger) AddPlannedCashWithdrawal(accountId uuid.UUID, withdrawal *PlannedCashFlow) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	if err := account.addPlannedCashWithdrawal(withdrawal); err != nil {
		return err
	}

	event := PlannedCashWithdrawalCreatedEvent{}.NewFrom(withdrawal)

	return l.addAndDispatchAccountEvent(account, event)
}

// LoadAccount an account by id
//
//
func (l *DefaultLedger) LoadAccount(accountId uuid.UUID) (*Account, error) {
	acc, err := l.accountRepository.loadById(accountId)
	if err != nil {
		return nil, err
	}

	if acc == nil {
		return nil, &AccountNotFoundError{accountId}
	}

	return acc, nil
}

func (l *DefaultLedger) addAndDispatchAccountEvent(account *Account, event events.Event) error {
	account.addRecordedEvent(event)

	if err := l.accountRepository.save(account); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(event)

	return nil
}

// ConfirmPlannedCashWithdrawal confirm a planned cash withdrawal
//
//
func (l *DefaultLedger) ConfirmPlannedCashWithdrawal(accountId uuid.UUID, withdrawalId uuid.UUID) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	flows := account.GetPlannedCashWithdrawals()
	withdrawal, ok := flows[withdrawalId]
	if !ok {
		return &PlannedCashWithdrawalNotFoundError{withdrawalId, accountId}
	}

	if err := account.addValue(withdrawal.amount, withdrawal.title); err != nil {
		return err
	}

	delete(account.plannedCashWithdrawals, withdrawalId)

	event := PlannedCashWithdrawalConfirmedEvent{}.NewFrom(withdrawal)

	if err := l.addAndDispatchAccountEvent(account, event); err != nil {
		return err
	}

	return nil
}
