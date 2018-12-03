package accounting

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/events"
	"testing"
)

var currentEvent events.Event

type TestEventHandler struct {
}

func (h *TestEventHandler) Handle(event events.Event) {
	currentEvent = event
}

// TestLedger_CreateAccount
//
//
func TestLedger_CreateAccount(t *testing.T) {
	eventDispatcher := events.DomainDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher, nil)

	eventDispatcher.RegisterHandler((&AccountCreatedEvent{}).GetName(), &TestEventHandler{})

	acc := ledger.CreateAccount("Yet another Bitcoin account", "BTC")
	assert.True(t, len(acc.id) == 16)
	assert.Equal(t, acc.title, "Yet another Bitcoin account")

	// event validation
	assert.IsType(t, &AccountCreatedEvent{}, currentEvent)
	assert.Equal(t, acc.balance, Money{}.NewFromInt(0, "BTC"))
	assert.Equal(t, acc.id, currentEvent.(*AccountCreatedEvent).accountId)
	assert.Equal(t, acc.title, currentEvent.(*AccountCreatedEvent).accountTitle)
	assert.Equal(t, "BTC", currentEvent.(*AccountCreatedEvent).currencyId)
}

// TestLedger_TransferValue
//
//
func TestLedger_TransferValue(t *testing.T) {
	eventDispatcher := events.DomainDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher, nil)

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
	eventDispatcher := events.DomainDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher, nil)

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
	eventDispatcher := events.DomainDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher, nil)

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

// TestLedger_LoadAccount
//
//
func TestLedger_LoadAccount(t *testing.T) {
	storage := &inMemoryEventStorage{}
	ledger := Ledger{}.New(events.DomainDispatcher{}.GetInstance(), storage)

	accountId := uuid.NewV4()

	// negative test when first event is not AccountCreatedEvent
	storage.AddEvent(&AccountValueAddedEvent{
		accountId: accountId,
		value:     Money{}.NewFromInt(1000000, "EUR"), // 10,000.00 EUR
		reason:    "initial",
	})
	assert.Len(t, ledger.eventStorage.(*inMemoryEventStorage).events, 1)

	_, err := ledger.LoadAccount(accountId)
	assert.IsType(t, &AccountCreatedEventNotFoundError{}, err)

	// positive test
	storage.events = []events.Event{} // clear old events

	storage.AddEvent(&AccountCreatedEvent{
		accountId:    accountId,
		accountTitle: "Test Account",
		currencyId:   "EUR",
	})
	storage.AddEvent(&AccountValueAddedEvent{
		accountId: accountId,
		value:     Money{}.NewFromInt(1000000, "EUR"), // 10,000.00 EUR
		reason:    "initial",
	})
	storage.AddEvent(&AccountValueSubtractedEvent{
		accountId: accountId,
		value:     Money{}.NewFromInt(90000, "EUR"), // 9,00.00 EUR (new account balance: 9,100.00 EUR)
		reason:    "monthly apartment rent",
	})
	storage.AddEvent(&AccountValueAddedEvent{
		accountId: accountId,
		value:     Money{}.NewFromInt(660000, "EUR"), // 6,600.00 EUR (new account balance: 15,700.00 EUR)
		reason:    "monthly salary",
	})
	storage.AddEvent(&AccountValueTransferredEvent{
		fromId: accountId,
		toId:   uuid.NewV4(),
		value:  Money{}.NewFromInt(100000, "EUR"), // 1,000.00 EUR (new account balance: 14,700.00 EUR)
		reason: "reserves",
	})
	storage.AddEvent(&AccountValueTransferredEvent{
		fromId: uuid.NewV4(),
		toId:   accountId,
		value:  Money{}.NewFromInt(50000, "EUR"), // 500.00 EUR (new account balance: 15,200.00 EUR)
		reason: "holidays",
	})

	// This two events should be ignored for our account, because they refer different accounts.
	storage.AddEvent(&AccountValueSubtractedEvent{
		accountId: uuid.NewV4(),
		value:     Money{}.NewFromInt(10000, "EUR"), // 100.00 EUR
		reason:    "what ever",
	})
	storage.AddEvent(&AccountValueTransferredEvent{
		fromId: uuid.NewV4(),
		toId:   uuid.NewV4(),
		value:  Money{}.NewFromInt(5000, "EUR"), // 50.00 EUR
		reason: "birthday",
	})

	assert.Len(t, ledger.eventStorage.(*inMemoryEventStorage).events, 8)

	// test history for account
	history := []events.Event{}

	for event := range ledger.getHistoryFor(accountId) {
		history = append(history, event)
	}
	assert.Len(t, history, 6)

	// try to load
	account, err := ledger.LoadAccount(accountId)
	assert.Nil(t, err)
	assert.Equal(t, accountId, account.id)
	assert.Equal(t, "Test Account", account.title)
	fmt.Printf("%s\n", account.balance.amount.String())
	assert.Equal(t, Money{}.NewFromInt(1520000, "EUR"), account.balance)
}
