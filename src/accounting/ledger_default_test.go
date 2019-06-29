package accounting

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/src/accounting/errors"
	events2 "pas/src/accounting/events"
	"pas/src/accounting/structs"
	"pas/src/events"
	"pas/src/money"
	errors2 "pas/src/money/errors"
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

// TestDefaultLedger_CreateAccount
//
//
func TestDefaultLedger_CreateAccount(t *testing.T) {
	eventDispatcher := events.DomainDispatcher{}.New()
	ledger := DefaultLedger{}.New(eventDispatcher, &TestEventStorage{})

	eventDispatcher.RegisterHandler((&events2.AccountCreatedEvent{}).GetName(), &TestEventHandler{})

	acc, _ := ledger.CreateAccount("Yet another Bitcoin account", "BTC")
	assert.True(t, len(acc.id) == 16)
	assert.Equal(t, acc.title, "Yet another Bitcoin account")

	// event validation
	assert.IsType(t, &events2.AccountCreatedEvent{}, currentEvent)
	assert.Equal(t, acc.balance, money.Money{}.NewFromInt(0, "BTC"))
	assert.Equal(t, acc.id, currentEvent.(*events2.AccountCreatedEvent).AccountId)
	assert.Equal(t, acc.title, currentEvent.(*events2.AccountCreatedEvent).AccountTitle)
	assert.Equal(t, "BTC", currentEvent.(*events2.AccountCreatedEvent).AurrencyId)
}

// TestDefaultLedger_TransferValue
//
//
func TestDefaultLedger_TransferValue(t *testing.T) {
	// prepare ledger and dispatcher
	eventDispatcher := events.DomainDispatcher{}.New()
	eventDispatcher.RegisterHandler((&events2.AccountValueTransferredEvent{}).GetName(), &TestEventHandler{})

	ledger := DefaultLedger{}.New(eventDispatcher, &TestEventStorage{})

	// our test accounts
	var fromAccount *Account
	var toAccount *Account

	// prepare accounts
	{
		from, err := ledger.CreateAccount("From Account", "EUR")
		to, err := ledger.CreateAccount("To Account", "EUR")
		assert.Nil(t, err)

		err = ledger.AddValue(from.GetId(), money.Money{}.NewFromInt(100000, "EUR"), "1000 EUR")
		err = ledger.AddValue(to.GetId(), money.Money{}.NewFromInt(500000, "EUR"), "5000 EUR")
		assert.Nil(t, err)

		fromAccount, _ = ledger.LoadAccount(from.GetId())
		toAccount, _ = ledger.LoadAccount(to.GetId())
	}

	// account not found error (fromAccount)
	{
		err := ledger.TransferValue(uuid.NewV4(), toAccount.GetId(), money.Money{}.NewFromInt(1, "EUR"), "")
		assert.IsType(t, &errors.AccountNotFoundError{}, err)
	}

	// account not found error (toAccount)
	{
		err := ledger.TransferValue(fromAccount.GetId(), uuid.NewV4(), money.Money{}.NewFromInt(1, "EUR"), "")
		assert.IsType(t, &errors.AccountNotFoundError{}, err)
	}

	// InsufficientFoundsError
	{
		value := money.Money{}.NewFromInt(100100, "EUR")
		err := ledger.TransferValue(fromAccount.GetId(), toAccount.GetId(), value, "foobar")
		assert.IsType(t, &errors.InsufficientFoundsError{}, err)
	}

	// UnequalCurrenciesError
	{
		value := money.Money{}.NewFromInt(99999, "USD")
		err := ledger.TransferValue(fromAccount.GetId(), toAccount.GetId(), value, "foobar")
		assert.IsType(t, &errors2.UnequalCurrenciesError{}, err)
	}

	// valid transfer
	{
		value := money.Money{}.NewFromInt(99999, "EUR")
		err := ledger.TransferValue(fromAccount.GetId(), toAccount.GetId(), value, "foobar")
		assert.Nil(t, err)

		// check new balances
		{
			from, _ := ledger.LoadAccount(fromAccount.GetId())
			to, _ := ledger.LoadAccount(toAccount.GetId())

			expectedFrom := money.Money{}.NewFromInt(1, "EUR")
			expectedTo := money.Money{}.NewFromInt(599999, "EUR")

			assert.True(t, from.GetBalance().IsEqual(expectedFrom))
			assert.True(t, to.GetBalance().IsEqual(expectedTo))

			// check dispatched event
			assert.IsType(t, &events2.AccountValueTransferredEvent{}, currentEvent)
		}
	}
}

// TestDefaultLedger_AddValue
//
//
func TestDefaultLedger_AddValue(t *testing.T) {
	eventDispatcher := events.DomainDispatcher{}.New()
	ledger := DefaultLedger{}.New(eventDispatcher, &TestEventStorage{})

	eventDispatcher.RegisterHandler((&events2.AccountValueAddedEvent{}).GetName(), &TestEventHandler{})

	acc, _ := ledger.CreateAccount("test account", "USD")

	// negative test with unequal currencies
	{
		wrongValue := money.Money{}.NewFromInt(1000, "EUR") // 10.00 EUR

		err := ledger.AddValue(acc.GetId(), wrongValue, "yehaaa")
		assert.IsType(t, &errors2.UnequalCurrenciesError{}, err)
	}

	// account not found test
	{
		_, err := ledger.LoadAccount(uuid.NewV4())
		assert.IsType(t, &errors.AccountNotFoundError{}, err)
	}

	// positive test #1
	{
		goodValue := money.Money{}.NewFromInt(500, "USD") // 5.00 USD

		err := ledger.AddValue(acc.GetId(), goodValue, "second try")
		assert.Nil(t, err)

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.True(t, acc.GetBalance().IsEqual(goodValue))

		// event validation
		assert.IsType(t, &events2.AccountValueAddedEvent{}, currentEvent)
		assert.Equal(t, acc.id, currentEvent.(*events2.AccountValueAddedEvent).AccountId)
		assert.Equal(t, goodValue, currentEvent.(*events2.AccountValueAddedEvent).Value)
		assert.Equal(t, "second try", currentEvent.(*events2.AccountValueAddedEvent).Reason)
	}

	// positive test #2
	{
		nextGoodValue := money.Money{}.NewFromInt(1000, "USD") // 10.00 USD

		err := ledger.AddValue(acc.GetId(), nextGoodValue, "third try")
		assert.Nil(t, err)

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.True(t, acc.GetBalance().IsEqual(money.Money{}.NewFromInt(1500, "USD")))

		// event validation
		assert.IsType(t, &events2.AccountValueAddedEvent{}, currentEvent)
		assert.Equal(t, acc.id, currentEvent.(*events2.AccountValueAddedEvent).AccountId)
		assert.Equal(t, nextGoodValue, currentEvent.(*events2.AccountValueAddedEvent).Value)
		assert.Equal(t, "third try", currentEvent.(*events2.AccountValueAddedEvent).Reason)
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

// TestDefaultLedger_SubtractValue
//
//
func TestDefaultLedger_SubtractValue(t *testing.T) {
	eventDispatcher := events.DomainDispatcher{}.New()
	ledger := DefaultLedger{}.New(eventDispatcher, &TestEventStorage{})

	eventDispatcher.RegisterHandler((&events2.AccountValueSubtractedEvent{}).GetName(), &TestEventHandler{})

	acc, _ := ledger.CreateAccount("test account", "USD")

	// negative test #1 (unequal currencies)
	{
		err := ledger.SubtractValue(acc.GetId(), money.Money{}.NewFromInt(1000, "EUR"), "just for fun") // 10.00 EUR
		assert.IsType(t, &errors2.UnequalCurrenciesError{}, err)
	}

	// negative test #2 (insufficient founds)
	{
		err := ledger.SubtractValue(acc.GetId(), money.Money{}.NewFromInt(1, "USD"), "test it") // 0.01 USD
		assert.IsType(t, &errors.InsufficientFoundsError{}, err)
	}

	// positive test
	{
		err := ledger.AddValue(acc.GetId(), money.Money{}.NewFromInt(10000, "USD"), "initial") // 100.00 USD
		assert.Nil(t, err)

		acc, _ := ledger.LoadAccount(acc.GetId())

		subValue := money.Money{}.NewFromInt(999, "USD")
		err = ledger.SubtractValue(acc.GetId(), subValue, "what ever") // 9.99 USD
		assert.Nil(t, err)

		expectedValue := money.Money{}.NewFromInt(9001, "USD")

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.Equal(t, expectedValue, acc.balance)

		// event validation
		assert.IsType(t, &events2.AccountValueSubtractedEvent{}, currentEvent)
		assert.Equal(t, acc.id, currentEvent.(*events2.AccountValueSubtractedEvent).AccountId)
		assert.Equal(t, subValue, currentEvent.(*events2.AccountValueSubtractedEvent).Value)
		assert.Equal(t, "what ever", currentEvent.(*events2.AccountValueSubtractedEvent).Reason)
	}
}

// TestDefaultLedger_LoadAccount
//
//
func TestDefaultLedger_LoadAccount(t *testing.T) {
	ledger := DefaultLedger{}.New(events.DomainDispatcher{}.New(), &TestEventStorage{})
	defaultLedger := ledger.(*DefaultLedger)
	storage := defaultLedger.accountRepository.eventStorage
	accountId := uuid.NewV4()

	// negative test when first event is not AccountCreatedEvent
	storage.AddEvent(&events2.AccountValueAddedEvent{
		AccountId: accountId,
		Value:     money.Money{}.NewFromInt(1000000, "EUR"), // 10,000.00 EUR
		Reason:    "initial",
	})
	assert.Len(t, storage.(*TestEventStorage).events, 1)

	_, err := defaultLedger.LoadAccount(accountId)
	assert.IsType(t, &errors.AccountCreatedEventNotFoundError{}, err)

	// positive test
	storage.(*TestEventStorage).events = []events.Event{} // clear old events

	storage.AddEvent(&events2.AccountCreatedEvent{
		AccountId:    accountId,
		AccountTitle: "Test Account",
		AurrencyId:   "EUR",
	})
	storage.AddEvent(&events2.AccountValueAddedEvent{
		AccountId: accountId,
		Value:     money.Money{}.NewFromInt(1000000, "EUR"), // 10,000.00 EUR
		Reason:    "initial",
	})
	storage.AddEvent(&events2.AccountValueSubtractedEvent{
		AccountId: accountId,
		Value:     money.Money{}.NewFromInt(90000, "EUR"), // 9,00.00 EUR (new account balance: 9,100.00 EUR)
		Reason:    "monthly apartment rent",
	})
	storage.AddEvent(&events2.AccountValueAddedEvent{
		AccountId: accountId,
		Value:     money.Money{}.NewFromInt(660000, "EUR"), // 6,600.00 EUR (new account balance: 15,700.00 EUR)
		Reason:    "monthly salary",
	})
	storage.AddEvent(&events2.AccountValueTransferredEvent{
		FromId: accountId,
		ToId:   uuid.NewV4(),
		Value:  money.Money{}.NewFromInt(100000, "EUR"), // 1,000.00 EUR (new account balance: 14,700.00 EUR)
		Reason: "reserves",
	})
	storage.AddEvent(&events2.AccountValueTransferredEvent{
		FromId: uuid.NewV4(),
		ToId:   accountId,
		Value:  money.Money{}.NewFromInt(50000, "EUR"), // 500.00 EUR (new account balance: 15,200.00 EUR)
		Reason: "holidays",
	})

	// This two events should be ignored for our account, because they refer different accounts.
	storage.AddEvent(&events2.AccountValueSubtractedEvent{
		AccountId: uuid.NewV4(),
		Value:     money.Money{}.NewFromInt(10000, "EUR"), // 100.00 EUR
		Reason:    "what ever",
	})
	storage.AddEvent(&events2.AccountValueTransferredEvent{
		FromId: uuid.NewV4(),
		ToId:   uuid.NewV4(),
		Value:  money.Money{}.NewFromInt(5000, "EUR"), // 50.00 EUR
		Reason: "birthday",
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
	fmt.Printf("%s\n", account.balance.GetAmount().String())
	assert.Equal(t, money.Money{}.NewFromInt(1520000, "EUR"), account.balance)
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
		amount := money.Money{}.NewFromInt(100000, "EUR") // 1,000.00 EUR
		receipt := structs.PlannedCashFlow{}.New(acc.GetId(), date, amount, title)

		err := ledger.AddPlannedCashReceipt(acc.GetId(), receipt)
		assert.IsType(t, &errors2.UnequalCurrenciesError{}, err)
	}

	// positive test with correct currency
	{
		amount := money.Money{}.NewFromInt(100000000, "BTC") // 1 BTC
		receipt := structs.PlannedCashFlow{}.New(acc.GetId(), date, amount, title)

		err := ledger.AddPlannedCashReceipt(acc.GetId(), receipt)
		assert.Nil(t, err)

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.Len(t, acc.plannedCashReceipts, 1)

		accReceipt := acc.plannedCashReceipts[receipt.GetId()]
		assert.Equal(t, date, accReceipt.Date)
		assert.Equal(t, title, accReceipt.Title)
		assert.Equal(t, amount, accReceipt.Amount)
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
		amount := money.Money{}.NewFromInt(100000, "EUR") // 1,000.00 EUR
		withdrawal := structs.PlannedCashFlow{}.New(acc.GetId(), date, amount, title)

		err := ledger.AddPlannedCashWithdrawal(acc.GetId(), withdrawal)
		assert.IsType(t, &errors2.UnequalCurrenciesError{}, err)
	}

	// positive test with correct currency
	{
		amount := money.Money{}.NewFromInt(100000000, "BTC") // 1 BTC
		withdrawal := structs.PlannedCashFlow{}.New(acc.GetId(), date, amount, title)

		err := ledger.AddPlannedCashWithdrawal(acc.GetId(), withdrawal)
		assert.Nil(t, err)

		acc, _ = ledger.LoadAccount(acc.GetId())
		assert.Len(t, acc.plannedCashWithdrawals, 1)

		accWithdrawal := acc.plannedCashWithdrawals[withdrawal.GetId()]
		assert.Equal(t, date, accWithdrawal.Date)
		assert.Equal(t, title, accWithdrawal.Title)
		assert.Equal(t, amount, accWithdrawal.Amount)
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

// TestDefaultLedger_ConfirmPlannedCashReceipt
//
//
func TestDefaultLedger_ConfirmPlannedCashReceipt(t *testing.T) {
	// prepare defaultLedger
	ledger := DefaultLedger{}.New(events.DomainDispatcher{}.New(), &TestEventStorage{})
	defaultLedger := ledger.(*DefaultLedger)

	// prepare test account
	acc, _ := defaultLedger.CreateAccount("Test account", "EUR")
	assert.Nil(t, acc.plannedCashWithdrawals)

	// add some initial balance
	initialBalance := money.Money{}.NewFromInt(10000, "EUR")
	err := defaultLedger.AddValue(acc.GetId(), initialBalance, "initial")
	assert.Nil(t, err)

	// test PlannedCashReceiptNotFoundError
	{
		err := defaultLedger.ConfirmPlannedCashReceipt(acc.GetId(), uuid.NewV4())
		assert.IsType(t, &errors.PlannedCashReceiptNotFoundError{}, err)
	}

	// test confirmation
	{
		// add planned receipt
		date := (time.Now()).Add(24 * time.Hour)
		amount := money.Money{}.NewFromInt(100000, "EUR") // 1,000.00 EUR
		plannedReceipt := structs.PlannedCashFlow{}.New(acc.GetId(), date, amount, "test receipt")

		err := defaultLedger.AddPlannedCashReceipt(acc.GetId(), plannedReceipt)
		assert.Nil(t, err)

		// reload account
		acc, err = defaultLedger.LoadAccount(acc.GetId())
		assert.Nil(t, err)
		assert.NotNil(t, acc)

		// check if initial balance did not changed already
		assert.True(t, acc.GetBalance().IsEqual(initialBalance))

		// confirm planned receipt
		err = ledger.ConfirmPlannedCashReceipt(acc.GetId(), plannedReceipt.GetId())
		assert.Nil(t, err)

		// reload account
		acc, err = defaultLedger.LoadAccount(acc.GetId())
		assert.Nil(t, err)
		assert.NotNil(t, acc)

		// check new balance
		expectedBalance := money.Money{}.NewFromInt(110000, "EUR")
		assert.True(t, acc.GetBalance().IsEqual(expectedBalance))
	}
}

// TestDefaultLedger_ConfirmPlannedCashWithdrawal
//
//
func TestDefaultLedger_ConfirmPlannedCashWithdrawal(t *testing.T) {
	// prepare defaultLedger
	ledger := DefaultLedger{}.New(events.DomainDispatcher{}.New(), &TestEventStorage{})
	defaultLedger := ledger.(*DefaultLedger)

	// prepare test account
	acc, _ := defaultLedger.CreateAccount("Test account", "EUR")
	assert.Nil(t, acc.plannedCashWithdrawals)

	// add some initial balance
	initialBalance := money.Money{}.NewFromInt(100000000, "EUR")
	err := defaultLedger.AddValue(acc.GetId(), initialBalance, "initial")
	assert.Nil(t, err)

	// test PlannedCashWithdrawalNotFoundError
	{
		err := defaultLedger.ConfirmPlannedCashWithdrawal(acc.GetId(), uuid.NewV4())
		assert.IsType(t, &errors.PlannedCashWithdrawalNotFoundError{}, err)
	}

	// test confirmation
	{
		// add planned withdrawal
		date := (time.Now()).Add(24 * time.Hour)
		amount := money.Money{}.NewFromInt(100000, "EUR") // 1,000.00 EUR
		plannedWithdrawal := structs.PlannedCashFlow{}.New(acc.GetId(), date, amount, "test withdrawal")

		err := defaultLedger.AddPlannedCashWithdrawal(acc.GetId(), plannedWithdrawal)
		assert.Nil(t, err)

		// reload account
		acc, err = defaultLedger.LoadAccount(acc.GetId())
		assert.Nil(t, err)
		assert.NotNil(t, acc)

		// check if initial balance did not changed already
		assert.True(t, acc.GetBalance().IsEqual(initialBalance))

		// confirm planned withdrawal
		err = ledger.ConfirmPlannedCashWithdrawal(acc.GetId(), plannedWithdrawal.GetId())
		assert.Nil(t, err)

		// reload account
		acc, err = defaultLedger.LoadAccount(acc.GetId())
		assert.Nil(t, err)
		assert.NotNil(t, acc)

		// check new balance
		expectedBalance := money.Money{}.NewFromInt(99900000, "EUR")
		assert.True(t, acc.GetBalance().IsEqual(expectedBalance))
	}
}
