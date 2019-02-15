package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
	"pas/accounting/errors"
	events2 "pas/accounting/events"
	"pas/accounting/structs"
	"pas/events"
	"pas/money"
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
		balance: money.Money{}.NewFromInt(0, currencyId),
	}

	event := &events2.AccountCreatedEvent{
		AccountId:    a.id,
		AccountTitle: title,
		AurrencyId:   currencyId,
	}

	return a, l.addAndDispatchAccountEvent(a, event)
}

// TransferValue transfer Value from one Account to another
//
//
func (l *DefaultLedger) TransferValue(fromAccountId, toAccountId uuid.UUID, value money.Money, reason string) error {
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
		return &errors.InsufficientFoundsError{}
	}

	currencyId := fromAccount.balance.GetCurrencyId()
	oldFromAmount := fromAccount.balance.GetAmount()
	oldToAmount := toAccount.balance.GetAmount()

	newFromAmount := (&big.Int{}).Sub(oldFromAmount, value.GetAmount())
	newToAmount := (&big.Int{}).Add(oldToAmount, value.GetAmount())

	newFromBalance, err := money.Money{}.NewFromString(newFromAmount.String(), currencyId)
	if err != nil {
		return err
	}
	newToBalance, err := money.Money{}.NewFromString(newToAmount.String(), currencyId)
	if err != nil {
		return err
	}

	fromAccount.balance = newFromBalance
	toAccount.balance = newToBalance

	event := &events2.AccountValueTransferredEvent{
		FromId: fromAccount.id,
		ToId:   toAccount.id,
		Value:  value,
		Reason: reason,
	}

	return l.addAndDispatchAccountEvent(fromAccount, event)
}

// AddValue add new Value to an account and dispatch AccountValueAddedEvent
//
//
func (l *DefaultLedger) AddValue(accountId uuid.UUID, value money.Money, reason string) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	if err := account.addValue(value, reason); err != nil {
		return err
	}

	event := &events2.AccountValueAddedEvent{
		AccountId: accountId,
		Value:     value,
		Reason:    reason,
	}

	return l.addAndDispatchAccountEvent(account, event)
}

// SubtractValue subtract Value from an account and dispatch AccountValueSubtractedEvent
//
//
func (l *DefaultLedger) SubtractValue(accountId uuid.UUID, value money.Money, reason string) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	if err := account.subtractValue(value, reason); err != nil {
		return err
	}

	event := &events2.AccountValueSubtractedEvent{
		AccountId: accountId,
		Value:     value,
		Reason:    reason,
	}

	return l.addAndDispatchAccountEvent(account, event)
}

// AddPlannedCashReceipt add a planned cash receipt to an account
//
//
func (l *DefaultLedger) AddPlannedCashReceipt(accountId uuid.UUID, receipt *structs.PlannedCashFlow) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	if err := account.addPlannedCashReceipt(receipt); err != nil {
		return err
	}

	event := events2.PlannedCashReceiptCreatedEvent{}.NewFrom(receipt)

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
		return &errors.PlannedCashReceiptNotFoundError{receiptId, accountId}
	}

	if err := account.addValue(receipt.Amount, receipt.Title); err != nil {
		return err
	}

	delete(account.plannedCashReceipts, receiptId)

	event := events2.PlannedCashReceiptConfirmedEvent{}.NewFrom(receipt)

	if err := l.addAndDispatchAccountEvent(account, event); err != nil {
		return err
	}

	return nil
}

// AddPlannedCashWithdrawal add a planned cash withdrawal to an account
//
//
func (l *DefaultLedger) AddPlannedCashWithdrawal(accountId uuid.UUID, withdrawal *structs.PlannedCashFlow) error {
	account, err := l.LoadAccount(accountId)
	if err != nil {
		return err
	}

	if err := account.addPlannedCashWithdrawal(withdrawal); err != nil {
		return err
	}

	event := events2.PlannedCashWithdrawalCreatedEvent{}.NewFrom(withdrawal)

	return l.addAndDispatchAccountEvent(account, event)
}

// LoadAccount an account by Id
//
//
func (l *DefaultLedger) LoadAccount(accountId uuid.UUID) (*Account, error) {
	acc, err := l.accountRepository.loadById(accountId)
	if err != nil {
		return nil, err
	}

	if acc == nil {
		return nil, &errors.AccountNotFoundError{accountId}
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
		return &errors.PlannedCashWithdrawalNotFoundError{withdrawalId, accountId}
	}

	if err := account.addValue(withdrawal.Amount, withdrawal.Title); err != nil {
		return err
	}

	delete(account.plannedCashWithdrawals, withdrawalId)

	event := events2.PlannedCashWithdrawalConfirmedEvent{}.NewFrom(withdrawal)

	if err := l.addAndDispatchAccountEvent(account, event); err != nil {
		return err
	}

	return nil
}
