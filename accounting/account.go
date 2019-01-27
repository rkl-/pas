package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
	"strings"
)

// Account a ledger accountId
//
//
type Account struct {
	id      uuid.UUID
	title   string
	balance Money
}

func (a *Account) GetId() uuid.UUID {
	return a.id
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
