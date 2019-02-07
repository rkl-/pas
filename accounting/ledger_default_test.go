package accounting

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/events"
	"testing"
	"time"
)

var currentEvent events.Event

type TestEventHandler struct {
}

func (h *TestEventHandler) Handle(event events.Event) {
	currentEvent = event
}

// TestEventStorage
//
//
type TestEventStorage struct {
	events []events.Event
}

func (s *TestEventStorage) AddEvent(event events.Event) {
	if s.events == nil {
		s.events = []events.Event{}
	}

	s.events = append(s.events, event)
}

func (s *TestEventStorage) GetEventStream() chan events.Event {
	ch := make(chan events.Event)

	go func() {
		defer close(ch)

		for _, event := range s.events {
			ch <- event
		}
	}()

	return ch
}

// TestLedger_CreateAccount
//
//
func TestLedger_CreateAccount(t *testing.T) {
	eventDispatcher := events.DomainDispatcher{}.New()
	ledger := DefaultLedger{}.New(eventDispatcher, &TestEventStorage{})

	eventDispatcher.RegisterHandler((&AccountCreatedEvent{}).GetName(), &TestEventHandler{})

	acc, _ := ledger.CreateAccount("Yet another Bitcoin account", "BTC")
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
	// prepare ledger and dispatcher
	eventDispatcher := events.DomainDispatcher{}.New()
	eventDispatcher.RegisterHandler((&AccountValueTransferredEvent{}).GetName(), &TestEventHandler{})

	ledger := DefaultLedger{}.New(eventDispatcher, &TestEventStorage{})

	// our test accounts
	var fromAccount *Account
	var toAccount *Account

	// prepare accounts
	{
		from, err := ledger.CreateAccount("From Account", "EUR")
		to, err := ledger.CreateAccount("To Account", "EUR")
		assert.Nil(t, err)

		err = ledger.AddValue(from.GetId(), Money{}.NewFromInt(100000, "EUR"), "1000 EUR")
		err = ledger.AddValue(to.GetId(), Money{}.NewFromInt(500000, "EUR"), "5000 EUR")
		assert.Nil(t, err)

		fromAccount, _ = ledger.LoadAccount(from.GetId())
		toAccount, _ = ledger.LoadAccount(to.GetId())
	}

	// account not found error (fromAccount)
	{
		err := ledger.TransferValue(uuid.NewV4(), toAccount.GetId(), Money{}.NewFromInt(1, "EUR"), "")
		assert.IsType(t, &AccountNotFoundError{}, err)
	}

	// account not found error (toAccount)
	{
		err := ledger.TransferValue(fromAccount.GetId(), uuid.NewV4(), Money{}.NewFromInt(1, "EUR"), "")
		assert.IsType(t, &AccountNotFoundError{}, err)
	}

	// InsufficientFoundsError
	{
		value := Money{}.NewFromInt(100100, "EUR")
		err := ledger.TransferValue(fromAccount.GetId(), toAccount.GetId(), value, "foobar")
		assert.IsType(t, &InsufficientFoundsError{}, err)
	}

	// UnequalCurrenciesError
	{
		value := Money{}.NewFromInt(99999, "USD")
		err := ledger.TransferValue(fromAccount.GetId(), toAccount.GetId(), value, "foobar")
		assert.IsType(t, &UnequalCurrenciesError{}, err)
	}

	// valid transfer
	{
		value := Money{}.NewFromInt(99999, "EUR")
		err := ledger.TransferValue(fromAccount.GetId(), toAccount.GetId(), value, "foobar")
		assert.Nil(t, err)

		// check new balances
		{
			from, _ := ledger.LoadAccount(fromAccount.GetId())
			to, _ := ledger.LoadAccount(toAccount.GetId())

			expectedFrom := Money{}.NewFromInt(1, "EUR")
			expectedTo := Money{}.NewFromInt(599999, "EUR")

			assert.True(t, from.GetBalance().IsEqual(expectedFrom))
			assert.True(t, to.GetBalance().IsEqual(expectedTo))

			// check dispatched event
			assert.IsType(t, &AccountValueTransferredEvent{}, currentEvent)
		}
	}
}

// TestLedger_AddValue
//
//
func TestLedger_AddValue(t *testing.T) {
	eventDispatcher := events.DomainDispatcher{}.New()
	ledger := DefaultLedger{}.New(eventDispatcher, &TestEventStorage{})

	eventDispatcher.RegisterHandler((&AccountValueAddedEvent{}).GetName(), &TestEventHandler{})

	acc, _ := ledger.CreateAccount("test account", "USD")

	// negative test with unequal currencies
	{
		wrongValue := Money{}.NewFromInt(1000, "EUR") // 10.00 EUR

		err := ledger.AddValue(acc.GetId(), wrongValue, "yehaaa")
		assert.IsType(t, &UnequalCurrenciesError{}, err)
	}

	// account not found test
	{
		_, err := ledger.LoadAccount(uuid.NewV4())
		assert.IsType(t, &AccountNotFoundError{}, err)
	}

	// positive test #1
	{
		goodValue := Money{}.NewFromInt(500, "USD") // 5.00 USD

		err := ledger.AddValue(acc.GetId(), goodValue, "second try")
		assert.Nil(t, err)

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.True(t, acc.GetBalance().IsEqual(goodValue))

		// event validation
		assert.IsType(t, &AccountValueAddedEvent{}, currentEvent)
		assert.Equal(t, acc.id, currentEvent.(*AccountValueAddedEvent).accountId)
		assert.Equal(t, goodValue, currentEvent.(*AccountValueAddedEvent).value)
		assert.Equal(t, "second try", currentEvent.(*AccountValueAddedEvent).reason)
	}

	// positive test #2
	{
		nextGoodValue := Money{}.NewFromInt(1000, "USD") // 10.00 USD

		err := ledger.AddValue(acc.GetId(), nextGoodValue, "third try")
		assert.Nil(t, err)

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.True(t, acc.GetBalance().IsEqual(Money{}.NewFromInt(1500, "USD")))

		// event validation
		assert.IsType(t, &AccountValueAddedEvent{}, currentEvent)
		assert.Equal(t, acc.id, currentEvent.(*AccountValueAddedEvent).accountId)
		assert.Equal(t, nextGoodValue, currentEvent.(*AccountValueAddedEvent).value)
		assert.Equal(t, "third try", currentEvent.(*AccountValueAddedEvent).reason)
	}

	// load account from ledger
	{
		reloadedAccount, err := ledger.LoadAccount(acc.GetId())
		assert.IsType(t, &Account{}, reloadedAccount)
		assert.Nil(t, err)
		assert.True(t, acc.balance.IsEqual(reloadedAccount.balance))
		assert.Equal(t, acc.title, reloadedAccount.title)
		assert.Equal(t, acc.plannedCashReceipts, reloadedAccount.plannedCashReceipts)
	}
}

// TestLedger_SubtractValue
//
//
func TestLedger_SubtractValue(t *testing.T) {
	eventDispatcher := events.DomainDispatcher{}.New()
	ledger := DefaultLedger{}.New(eventDispatcher, &TestEventStorage{})

	eventDispatcher.RegisterHandler((&AccountValueSubtractedEvent{}).GetName(), &TestEventHandler{})

	acc, _ := ledger.CreateAccount("test account", "USD")

	// negative test #1 (unequal currencies)
	{
		err := ledger.SubtractValue(acc.GetId(), Money{}.NewFromInt(1000, "EUR"), "just for fun") // 10.00 EUR
		assert.IsType(t, &UnequalCurrenciesError{}, err)
	}

	// negative test #2 (insufficient founds)
	{
		err := ledger.SubtractValue(acc.GetId(), Money{}.NewFromInt(1, "USD"), "test it") // 0.01 USD
		assert.IsType(t, &InsufficientFoundsError{}, err)
	}

	// positive test
	{
		err := ledger.AddValue(acc.GetId(), Money{}.NewFromInt(10000, "USD"), "initial") // 100.00 USD
		assert.Nil(t, err)

		acc, _ := ledger.LoadAccount(acc.GetId())

		subValue := Money{}.NewFromInt(999, "USD")
		err = ledger.SubtractValue(acc.GetId(), subValue, "what ever") // 9.99 USD
		assert.Nil(t, err)

		expectedValue := Money{}.NewFromInt(9001, "USD")

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.Equal(t, expectedValue, acc.balance)

		// event validation
		assert.IsType(t, &AccountValueSubtractedEvent{}, currentEvent)
		assert.Equal(t, acc.id, currentEvent.(*AccountValueSubtractedEvent).accountId)
		assert.Equal(t, subValue, currentEvent.(*AccountValueSubtractedEvent).value)
		assert.Equal(t, "what ever", currentEvent.(*AccountValueSubtractedEvent).reason)
	}
}

// TestLedger_LoadAccount
//
//
func TestLedger_LoadAccount(t *testing.T) {
	ledger := DefaultLedger{}.New(events.DomainDispatcher{}.New(), &TestEventStorage{})
	defaultLedger := ledger.(*DefaultLedger)
	storage := defaultLedger.accountRepository.eventStorage
	accountId := uuid.NewV4()

	// negative test when first event is not AccountCreatedEvent
	storage.AddEvent(&AccountValueAddedEvent{
		accountId: accountId,
		value:     Money{}.NewFromInt(1000000, "EUR"), // 10,000.00 EUR
		reason:    "initial",
	})
	assert.Len(t, storage.(*TestEventStorage).events, 1)

	_, err := defaultLedger.LoadAccount(accountId)
	assert.IsType(t, &AccountCreatedEventNotFoundError{}, err)

	// positive test
	storage.(*TestEventStorage).events = []events.Event{} // clear old events

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

	assert.Len(t, storage.(*TestEventStorage).events, 8)

	// test history for account
	history := []events.Event{}

	for event := range defaultLedger.accountRepository.getHistoryFor(accountId) {
		history = append(history, event)
	}
	assert.Len(t, history, 6)

	// try to load
	account, err := defaultLedger.LoadAccount(accountId)
	assert.Nil(t, err)
	assert.Equal(t, accountId, account.id)
	assert.Equal(t, "Test Account", account.title)
	fmt.Printf("%s\n", account.balance.amount.String())
	assert.Equal(t, Money{}.NewFromInt(1520000, "EUR"), account.balance)
}

// TestDefaultLedger_AddPlannedCashReceipt
//
//
func TestDefaultLedger_AddPlannedCashReceipt(t *testing.T) {
	// prepare defaultLedger
	ledger := DefaultLedger{}.New(events.DomainDispatcher{}.New(), &TestEventStorage{})
	defaultLedger := ledger.(*DefaultLedger)

	// prepare test account
	acc, _ := defaultLedger.CreateAccount("Test account", "BTC")
	assert.Nil(t, acc.plannedCashReceipts)

	// prepare some base details
	date := (time.Now()).Add(24 * time.Hour)
	title := "test receipt"

	// negative test with wrong currency
	{
		amount := Money{}.NewFromInt(100000, "EUR") // 1,000.00 EUR
		receipt := PlannedCashFlow{}.New(acc.GetId(), date, amount, title)

		err := ledger.AddPlannedCashReceipt(acc.GetId(), receipt)
		assert.IsType(t, &UnequalCurrenciesError{}, err)
	}

	// positive test with correct currency
	{
		amount := Money{}.NewFromInt(100000000, "BTC") // 1 BTC
		receipt := PlannedCashFlow{}.New(acc.GetId(), date, amount, title)

		err := ledger.AddPlannedCashReceipt(acc.GetId(), receipt)
		assert.Nil(t, err)

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.Len(t, acc.plannedCashReceipts, 1)

		accReceipt := acc.plannedCashReceipts[receipt.GetId()]
		assert.Equal(t, date, accReceipt.date)
		assert.Equal(t, title, accReceipt.title)
		assert.Equal(t, amount, accReceipt.amount)
	}

	// load account from ledger
	{
		reloadedAccount, err := ledger.LoadAccount(acc.GetId())
		assert.IsType(t, &Account{}, reloadedAccount)
		assert.Nil(t, err)
		assert.True(t, acc.balance.IsEqual(reloadedAccount.balance))
		assert.Equal(t, acc.title, reloadedAccount.title)
		assert.Equal(t, acc.plannedCashReceipts, reloadedAccount.plannedCashReceipts)
	}
}

// TestDefaultLedger_AddPlannedCashWithdrawal
//
//
func TestDefaultLedger_AddPlannedCashWithdrawal(t *testing.T) {
	// prepare defaultLedger
	ledger := DefaultLedger{}.New(events.DomainDispatcher{}.New(), &TestEventStorage{})
	defaultLedger := ledger.(*DefaultLedger)

	// prepare test account
	acc, _ := defaultLedger.CreateAccount("Test account", "BTC")
	assert.Nil(t, acc.plannedCashWithdrawals)

	// prepare some base details
	date := (time.Now()).Add(24 * time.Hour)
	title := "test withdrawal"

	// negative test with wrong currency
	{
		amount := Money{}.NewFromInt(100000, "EUR") // 1,000.00 EUR
		withdrawal := PlannedCashFlow{}.New(acc.GetId(), date, amount, title)

		err := ledger.AddPlannedCashWithdrawal(acc.GetId(), withdrawal)
		assert.IsType(t, &UnequalCurrenciesError{}, err)
	}

	// positive test with correct currency
	{
		amount := Money{}.NewFromInt(100000000, "BTC") // 1 BTC
		withdrawal := PlannedCashFlow{}.New(acc.GetId(), date, amount, title)

		err := ledger.AddPlannedCashWithdrawal(acc.GetId(), withdrawal)
		assert.Nil(t, err)

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.Len(t, acc.plannedCashWithdrawals, 1)

		accWithdrawal := acc.plannedCashWithdrawals[withdrawal.GetId()]
		assert.Equal(t, date, accWithdrawal.date)
		assert.Equal(t, title, accWithdrawal.title)
		assert.Equal(t, amount, accWithdrawal.amount)
	}

	// load account from ledger
	{
		reloadedAccount, err := ledger.LoadAccount(acc.GetId())
		assert.IsType(t, &Account{}, reloadedAccount)
		assert.Nil(t, err)
		assert.True(t, acc.balance.IsEqual(reloadedAccount.balance))
		assert.Equal(t, acc.title, reloadedAccount.title)
		assert.Equal(t, acc.plannedCashWithdrawals, reloadedAccount.plannedCashWithdrawals)
	}
}
