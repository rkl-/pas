package cq_command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCreateAccountCommand_New
//
//
func TestCreateAccountCommand_New(t *testing.T) {
	cmd := CreateAccountCommand{}.New("Test title", "BTC")
	assert.Equal(t, "Test title", cmd.title)
	assert.Equal(t, "BTC", cmd.currencyId)
}

// TestCreateAccountCommand_GetRequestId
//
//
func TestCreateAccountCommand_GetRequestId(t *testing.T) {
	cmd := CreateAccountCommand{}.New("Test title", "BTC")
	assert.Equal(t, "command.create_account", cmd.GetRequestId())
}
