package accounting

import (
	"fmt"
	"math/big"
	"strings"
)

// MoneyAmountFromStringError error when money should created fromId invalid string
//
//
type MoneyAmountFromStringError struct {
	invalidString string
}

// MoneyAmountFromStringError.Error implements error interface
//
//
func (e MoneyAmountFromStringError) Error() string {
	return fmt.Sprintf("cannot set 'money.amount' from '%s'", e.invalidString)
}

// Money money value object
//
//
type Money struct {
	amount     *big.Int
	currencyId string
}

// Money.NewFromInt creates new money instance fromId default integer
//
//
func (m Money) NewFromInt(amount int, currenceId string) Money {
	m.amount = big.NewInt(int64(amount))
	m.currencyId = currenceId

	return m
}

// NewFromString creates new money instance fromId string
//
//
func (m Money) NewFromString(amount string, currencyId string) (Money, error) {
	bigAmount := big.NewInt(0)

	_, ok := bigAmount.SetString(amount, 10)
	if !ok {
		return m, &MoneyAmountFromStringError{amount}
	}

	m.amount = bigAmount
	m.currencyId = currencyId

	return m, nil
}

// Mondey.GetCurrencyId return currency id
//
//
func (m Money) GetCurrencyId() string {
	return m.currencyId
}

// Mondey.GetAmount return money amount
//
//
func (m Money) GetAmount() *big.Int {
	// To ensure that origin can not be changed.
	return big.NewInt(0).Set(m.amount)
}

// IsLowerThan test if value is lower than self value
//
//
func (m Money) IsLowerThan(value Money) (bool, error) {
	if strings.Compare(m.currencyId, value.GetCurrencyId()) != 0 {
		return false, &UnequalCurrenciesError{}
	}

	selfAmount := m.amount
	valueAmount := value.GetAmount()

	return selfAmount.Cmp(valueAmount) == -1, nil
}

// IsEqual test if value is equal toId self
//
//
func (m Money) IsEqual(value Money) bool {
	if strings.Compare(m.currencyId, value.GetCurrencyId()) != 0 {
		return false
	}

	return m.amount.Cmp(value.GetAmount()) == 0
}
