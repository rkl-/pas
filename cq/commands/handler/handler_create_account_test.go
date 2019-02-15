package handler

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/accounting"
	"pas/cq"
	commandPkg "pas/cq/commands/command"
	"pas/events"
	"pas/money"
	"testing"
)

// TestCreateAccountCommandHandler_Handle
//
//
func TestCreateAccountCommandHandler_Handle(t *testing.T) {
	// Prepare ledger
	dispatcher := events.DomainDispatcher{}.New()
	ledger := accounting.DefaultLedger{}.New(dispatcher, &events.InMemoryEventStorage{})

	// Prepare command bus and register command handler
	cmdBus := cq.CommandBus{}.New()
	err := cmdBus.RegisterHandler("command.unsupported_command", &CreateAccountCommandHandler{
		ledger: ledger,
	})
	assert.Nil(t, err)

	err = cmdBus.RegisterHandler("command.create_account", &CreateAccountCommandHandler{
		ledger: ledger,
	})
	assert.Nil(t, err)

	// negative test for UnsupportedRequestError
	{
		_, err = cmdBus.Execute(&unsupportedCommand{})
		assert.IsType(t, &cq.UnsupportedRequestError{}, err)
	}

	// positive test
	{
		title := "My new test account"
		currencyId := "BTC"

		command := commandPkg.CreateAccountCommand{}.New(title, currencyId)

		accountId, err := cmdBus.Execute(command)
		assert.Nil(t, err)
		assert.IsType(t, uuid.UUID{}, accountId)

		// load created account from ledger and compare with expected one
		{
			loadedAccount, err := ledger.LoadAccount(accountId.(uuid.UUID))
			assert.Nil(t, err)
			assert.NotNil(t, loadedAccount)
			assert.Equal(t, title, loadedAccount.GetTitle())
			assert.Equal(t, currencyId, loadedAccount.GetCurrencyId())

			// It also should have a zero balance.
			assert.True(t, loadedAccount.GetBalance().IsEqual(money.Money{}.NewFromInt(0, currencyId)))
		}
	}
}
