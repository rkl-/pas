package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
	"pas/accounting/errors"
	"pas/accounting/structs"
	"pas/events"
	"pas/money"
	errors2 "pas/money/errors"
	"strings"
)

// Account a ledger account
//
//
type Account struct {
	id                     uuid.UUID
	title                  string
	balance                money.Money
	plannedCashReceipts    structs.PlannedCashFlowMap
	plannedCashWithdrawals structs.PlannedCashFlowMap
	recordedEvents         []events.Event
}

func (a *Account) GetId() uuid.UUID {
	return a.id
}

func (a *Account) GetTitle() string {
	return a.title
}

func (a *Account) GetCurrencyId() string {
	return a.balance.GetCurrencyId()
}

func (a *Account) GetBalance() money.Money {
	return a.balance
}

func (a *Account) GetPlannedCashReceipts() structs.PlannedCashFlowMap {
	return a.plannedCashReceipts
}

func (a *Account) GetPlannedCashWithdrawals() structs.PlannedCashFlowMap {
	return a.plannedCashWithdrawals
}

func (a *Account) addRecordedEvent(event events.Event) {
	if a.recordedEvents == nil {
		a.recordedEvents = []events.Event{}
	}

	a.recordedEvents = append(a.recordedEvents, event)
}

func (a *Account) addValue(value money.Money, reason string) error {
	if strings.Compare(a.balance.GetCurrencyId(), value.GetCurrencyId()) != 0 {
		return &errors2.UnequalCurrenciesError{}
	}

	currencyId := a.balance.GetCurrencyId()
	newAmount := (&big.Int{}).Add(a.balance.GetAmount(), value.GetAmount())

	newBalance, err := money.Money{}.NewFromString(newAmount.String(), currencyId)
	if err != nil {
		return err
	}

	a.balance = newBalance

	return nil
}

func (a *Account) subtractValue(value money.Money, reason string) error {
	ok, err := a.balance.IsLowerThan(value)
	if err != nil {
		return err
	}
	if ok {
		return &errors.InsufficientFoundsError{}
	}

	currencyId := a.balance.GetCurrencyId()
	newAmount := (&big.Int{}).Sub(a.balance.GetAmount(), value.GetAmount())

	newBalance, err := money.Money{}.NewFromString(newAmount.String(), currencyId)
	if err != nil {
		return err
	}

	a.balance = newBalance

	return nil
}

func (a *Account) addPlannedCashFlow(cashFlow *structs.PlannedCashFlow, target *structs.PlannedCashFlowMap) error {
	amount := cashFlow.Amount

	if strings.Compare(a.balance.GetCurrencyId(), amount.GetCurrencyId()) != 0 {
		return &errors2.UnequalCurrenciesError{}
	}

	if *target == nil {
		*target = structs.PlannedCashFlowMap{}
	}

	(*target)[cashFlow.GetId()] = cashFlow

	return nil
}

func (a *Account) addPlannedCashReceipt(receipt *structs.PlannedCashFlow) error {
	return a.addPlannedCashFlow(receipt, &a.plannedCashReceipts)
}

func (a *Account) addPlannedCashWithdrawal(withdrawal *structs.PlannedCashFlow) error {
	return a.addPlannedCashFlow(withdrawal, &a.plannedCashWithdrawals)
}
