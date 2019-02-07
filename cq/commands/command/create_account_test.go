package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCreateAccountCommand_New
//
//
func TestCreateAccountCommand_New(t *testing.T) {
	cmd := CreateAccountCommand{}.New("Test Title", "BTC")
	assert.Equal(t, "Test Title", cmd.Title)
	assert.Equal(t, "BTC", cmd.CurrencyId)
}

// TestCreateAccountCommand_GetRequestId
//
//
func TestCreateAccountCommand_GetRequestId(t *testing.T) {
	cmd := CreateAccountCommand{}.New("Test Title", "BTC")
	assert.Equal(t, "command.create_account", cmd.GetRequestId())
}
