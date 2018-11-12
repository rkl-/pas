package accounting

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestMoney_NewFromInt
//
//
func TestMoney_NewFromInt(t *testing.T) {
	testAmount := 999999
	testCurrency := "BTC"

	testMoney := Money{}.NewFromInt(testAmount, testCurrency)
	testMoneyAmount := testMoney.GetAmount()
	testMoneyAmount.Add(testMoneyAmount, testMoneyAmount)

	newTestMoneyAmount := testMoney.GetAmount()

	assert.Equal(t, testAmount, int(newTestMoneyAmount.Int64()))
	assert.Equal(t, testCurrency, testMoney.GetCurrencyId())
}

// TestMoney_NewFromString
//
//
func TestMoney_NewFromString(t *testing.T) {
	testAmount := "9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"
	testCurrency := "ETH"

	// success
	testMoney, err := Money{}.NewFromString(testAmount, testCurrency)
	assert.Nil(t, err)
	assert.Equal(t, testAmount, testMoney.GetAmount().String())
	assert.Equal(t, testCurrency, testMoney.GetCurrencyId())

	// error
	testInvalidAmount := "foo-bar-what-ever"
	testMoney, err = Money{}.NewFromString(testInvalidAmount, testCurrency)
	assert.IsType(t, &MoneyAmountFromStringError{}, err)
	assert.Equal(t, fmt.Sprintf("cannot set 'money.amount' from '%s'", testInvalidAmount), err.Error())
}
