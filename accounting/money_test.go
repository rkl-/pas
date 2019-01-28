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
	assert.Equal(t, fmt.Sprintf("cannot set 'money.Amount' from '%s'", testInvalidAmount), err.Error())
}

// TestMoney_IsLowerThan
//
//
func TestMoney_IsLowerThan(t *testing.T) {
	// positive test
	x := Money{}.NewFromInt(100, "EUR") // 1.00 EUR
	y := Money{}.NewFromInt(101, "EUR") // 1.01 EUR

	ok, _ := x.IsLowerThan(y)
	assert.True(t, ok)

	// negative test #1
	y = Money{}.NewFromInt(100, "EUR") // 1.00 EUR

	ok, _ = x.IsLowerThan(y)
	assert.False(t, ok)

	// negative test #2
	y = Money{}.NewFromInt(99, "EUR") // 0.99 EUR

	ok, _ = x.IsLowerThan(y)
	assert.False(t, ok)

	// currency mismatch test
	y = Money{}.NewFromInt(999, "USD") // 9.99 USD

	_, err := x.IsLowerThan(y)
	assert.IsType(t, &UnequalCurrenciesError{}, err)
}

// TestMoney_IsEqual
//
//
func TestMoney_IsEqual(t *testing.T) {
	// positive test
	x := Money{}.NewFromInt(1000, "EUR") // 10.00 EUR
	y := Money{}.NewFromInt(1000, "EUR") // 10.00 EUR
	assert.True(t, x.IsEqual(y))

	// negative test #1
	y = Money{}.NewFromInt(1000, "USD") // 10.00 USD
	assert.False(t, x.IsEqual(y))

	// negative test #2
	y = Money{}.NewFromInt(1001, "EUR") // 10.01 EUR
	assert.False(t, x.IsEqual(y))
}
