package accounting

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testLedgerCreateAccountHandlerExecuted = false

type TestAccountCreatedEventHandler struct {
}

func (h *TestAccountCreatedEventHandler) Handle(event EventInterface) {
	testLedgerCreateAccountHandlerExecuted = true
}

// TestLedger_CreateAccount
//
//
func TestLedger_CreateAccount(t *testing.T) {
	eventDispatcherInstance = nil

	eventDispatcher := EventDispatcher{}.GetInstance()
	ledger := Ledger{}.New(eventDispatcher)

	eventDispatcher.RegisterHandler((&AccountCreatedEvent{}).GetName(), &TestAccountCreatedEventHandler{})

	acc := ledger.CreateAccount("Yet another Bitcoin account", "BTC")
	assert.True(t, len(acc.id) == 16)
	assert.Equal(t, acc.title, "Yet another Bitcoin account")
	assert.Equal(t, acc.balance, Money{}.NewFromInt(0, "BTC"))
	assert.True(t, testLedgerCreateAccountHandlerExecuted)
}

// TestLedger_TransferValue
//
//
func TestLedger_TransferValue(t *testing.T) {
	ledger := Ledger{}.New(EventDispatcher{}.GetInstance())

	// positive test
	fromAccount := ledger.CreateAccount("from account", "EUR")
	fromAccount.balance = Money{}.NewFromInt(100000, "EUR") // 1000.00 EUR

	toAccount := ledger.CreateAccount("to account", "EUR")
	toAccount.balance = Money{}.NewFromInt(50000, "EUR") // 500.00 EUR

	err := ledger.TransferValue(fromAccount, toAccount, Money{}.NewFromInt(10000, "EUR"), "foobar") // 100.00 EUR
	assert.Nil(t, err)
	assert.True(t, fromAccount.balance.IsEqual(Money{}.NewFromInt(90000, "EUR"))) // 900.00 EUR
	assert.True(t, toAccount.balance.IsEqual(Money{}.NewFromInt(60000, "EUR")))   // 600.00 EUR

	// TODO, test dispatching of AccountValueTransferEvent
}
