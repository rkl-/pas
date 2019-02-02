package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
	"pas/events"
	"strings"
)

// Account a ledger account
//
//
type Account struct {
	id                  uuid.UUID
	title               string
	balance             Money
	plannedCashReceipts []*PlannedCashReceipt
	recordedEvents      []events.Event
}

func (a *Account) GetId() uuid.UUID {
	return a.id
}

func (a *Account) GetTitle() string {
	return a.title
}

func (a *Account) GetCurrencyId() string {
	return a.balance.currencyId
}

func (a *Account) GetBalance() Money {
	return a.balance
}

func (a *Account) addRecordedEvent(event events.Event) {
	if a.recordedEvents == nil {
		a.recordedEvents = []events.Event{}
	}

	a.recordedEvents = append(a.recordedEvents, event)
}

func (a *Account) addValue(value Money, reason string) error {
	if strings.Compare(a.balance.currencyId, value.currencyId) != 0 {
		return &UnequalCurrenciesError{}
	}

	currencyId := a.balance.currencyId
	newAmount := (&big.Int{}).Add(a.balance.amount, value.amount)

	newBalance, err := Money{}.NewFromString(newAmount.String(), currencyId)
	if err != nil {
		return err
	}

	a.balance = newBalance

	return nil
}

func (a *Account) subtractValue(value Money, reason string) error {
	ok, err := a.balance.IsLowerThan(value)
	if err != nil {
		return err
	}
	if ok {
		return &InsufficientFoundsError{}
	}

	currencyId := a.balance.currencyId
	newAmount := (&big.Int{}).Sub(a.balance.amount, value.amount)

	newBalance, err := Money{}.NewFromString(newAmount.String(), currencyId)
	if err != nil {
		return err
	}

	a.balance = newBalance

	return nil
}

func (a *Account) addPlannedCashReceipt(receipt *PlannedCashReceipt) error {
	amount := receipt.amount

	if strings.Compare(a.balance.currencyId, amount.currencyId) != 0 {
		return &UnequalCurrenciesError{}
	}

	if a.plannedCashReceipts == nil {
		a.plannedCashReceipts = []*PlannedCashReceipt{}
	}

	a.plannedCashReceipts = append(a.plannedCashReceipts, receipt)

	return nil
}
