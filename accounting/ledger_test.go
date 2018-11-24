package accounting

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var currentEvent EventInterface

type TestEventHandler struct {
}

func (h *TestEventHandler) Handle(event EventInterface) {
	currentEvent = event
}

// TestLedger_CreateAccount
//
//
func TestLedger_CreateAccount(t *testing.T) {
	eventDispatcherInstance = nil

	eventDispatcher := EventDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher)

	eventDispatcher.RegisterHandler((&AccountCreatedEvent{}).GetName(), &TestEventHandler{})

	acc := ledger.CreateAccount("Yet another Bitcoin account", "BTC")
	assert.True(t, len(acc.id) == 16)
	assert.Equal(t, acc.title, "Yet another Bitcoin account")

	// event validation
	assert.IsType(t, &AccountCreatedEvent{}, currentEvent)
	assert.Equal(t, acc.balance, Money{}.NewFromInt(0, "BTC"))
	assert.Equal(t, acc.id, currentEvent.(*AccountCreatedEvent).accountId)
	assert.Equal(t, acc.title, currentEvent.(*AccountCreatedEvent).accountTitle)
}

// TestLedger_TransferValue
//
//
func TestLedger_TransferValue(t *testing.T) {
	eventDispatcher := EventDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher)

	eventDispatcher.RegisterHandler((&AccountValueTransferredEvent{}).GetName(), &TestEventHandler{})

	// positive test
	fromAccount := ledger.CreateAccount("fromId account", "EUR")
	fromAccount.balance = Money{}.NewFromInt(100000, "EUR") // 1000.00 EUR

	toAccount := ledger.CreateAccount("toId account", "EUR")
	toAccount.balance = Money{}.NewFromInt(50000, "EUR") // 500.00 EUR

	transferValue := Money{}.NewFromInt(10000, "EUR")
	err := ledger.TransferValue(fromAccount, toAccount, transferValue, "foobar") // 100.00 EUR
	assert.Nil(t, err)
	assert.True(t, fromAccount.balance.IsEqual(Money{}.NewFromInt(90000, "EUR"))) // 900.00 EUR
	assert.True(t, toAccount.balance.IsEqual(Money{}.NewFromInt(60000, "EUR")))   // 600.00 EUR

	// event validation
	assert.IsType(t, &AccountValueTransferredEvent{}, currentEvent)
	assert.Equal(t, fromAccount.id, currentEvent.(*AccountValueTransferredEvent).fromId)
	assert.Equal(t, toAccount.id, currentEvent.(*AccountValueTransferredEvent).toId)
	assert.Equal(t, transferValue, currentEvent.(*AccountValueTransferredEvent).value)
	assert.Equal(t, "foobar", currentEvent.(*AccountValueTransferredEvent).reason)
}

// TestLedger_AddValue
//
//
func TestLedger_AddValue(t *testing.T) {
	eventDispatcher := EventDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher)

	eventDispatcher.RegisterHandler((&AccountValueAddedEvent{}).GetName(), &TestEventHandler{})

	// negative test with unequal currencies
	acc := ledger.CreateAccount("test account", "USD")
	wrongValue := Money{}.NewFromInt(1000, "EUR") // 10.00 EUR

	err := ledger.AddValue(acc, wrongValue, "yehaaa")
	assert.IsType(t, &UnequalCurrenciesError{}, err)

	// positive test #1
	goodValue := Money{}.NewFromInt(500, "USD") // 5.00 USD

	err = ledger.AddValue(acc, goodValue, "second try")
	assert.Nil(t, err)
	assert.Equal(t, acc.balance, goodValue)

	// event validation
	assert.IsType(t, &AccountValueAddedEvent{}, currentEvent)
	assert.Equal(t, acc.id, currentEvent.(*AccountValueAddedEvent).accountId)
	assert.Equal(t, goodValue, currentEvent.(*AccountValueAddedEvent).value)
	assert.Equal(t, "second try", currentEvent.(*AccountValueAddedEvent).reason)

	// positive test #2
	nextGoodValue := Money{}.NewFromInt(1000, "USD") // 10.00 USD

	err = ledger.AddValue(acc, nextGoodValue, "third try")
	assert.Nil(t, err)
	assert.Equal(t, acc.balance, Money{}.NewFromInt(1500, "USD"))

	// event validation
	assert.IsType(t, &AccountValueAddedEvent{}, currentEvent)
	assert.Equal(t, acc.id, currentEvent.(*AccountValueAddedEvent).accountId)
	assert.Equal(t, nextGoodValue, currentEvent.(*AccountValueAddedEvent).value)
	assert.Equal(t, "third try", currentEvent.(*AccountValueAddedEvent).reason)
}

// TestLedger_SubtractValue
//
//
func TestLedger_SubtractValue(t *testing.T) {
	eventDispatcher := EventDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher)

	eventDispatcher.RegisterHandler((&AccountValueSubtractedEvent{}).GetName(), &TestEventHandler{})

	// negative test #1 (unequal currencies)
	acc := ledger.CreateAccount("test account", "USD")
	err := ledger.SubtractValue(acc, Money{}.NewFromInt(1000, "EUR"), "just for fun") // 10.00 EUR
	assert.IsType(t, &UnequalCurrenciesError{}, err)

	// negative test #2 (insufficient founds)
	err = ledger.SubtractValue(acc, Money{}.NewFromInt(1, "USD"), "test it") // 0.01 USD
	assert.IsType(t, &InsufficientFoundsError{}, err)

	// positive test
	err = ledger.AddValue(acc, Money{}.NewFromInt(10000, "USD"), "initial") // 100.00 USD
	assert.Nil(t, err)

	subValue := Money{}.NewFromInt(999, "USD")
	err = ledger.SubtractValue(acc, subValue, "what ever") // 9.99 USD
	assert.Nil(t, err)

	expectedValue := Money{}.NewFromInt(9001, "USD")
	assert.Equal(t, expectedValue, acc.balance)

	// event validation
	assert.IsType(t, &AccountValueSubtractedEvent{}, currentEvent)
	assert.Equal(t, acc.id, currentEvent.(*AccountValueSubtractedEvent).accountId)
	assert.Equal(t, subValue, currentEvent.(*AccountValueSubtractedEvent).value)
	assert.Equal(t, "what ever", currentEvent.(*AccountValueSubtractedEvent).reason)
}
