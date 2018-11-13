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
