package accounting

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestLedger_CreateAccount
//
//
func TestLedger_CreateAccount(t *testing.T) {
	acc := (&Ledger{}).CreateAccount("Yet another Bitcoin account", "BTC")
	assert.True(t, len(acc.id) == 16)
	assert.Equal(t, acc.title, "Yet another Bitcoin account")
	assert.Equal(t, acc.balance, Money{}.NewFromInt(0, "BTC"))

	// TODO, test if AccountCreatedEvent was dispatched
}
