package handler

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/accounting"
	"pas/cq"
	"pas/cq/commands/command"
	"pas/events"
	"testing"
	"time"
)

// TestConfirmPlannedCashWithdrawalCommandHandler_Handle
//
//
func TestConfirmPlannedCashWithdrawalCommandHandler_Handle(t *testing.T) {
	// prepare ledger
	dispatcher := events.DomainDispatcher{}.New()
	ledger := accounting.DefaultLedger{}.New(dispatcher, &events.InMemoryEventStorage{})

	// negative test for UnsupportedRequestError
	{
		cmdBus := cq.CommandBus{}.New()
		err := cmdBus.RegisterHandler("command.unsupported_command", &CreatePlannedCashWithdrawalCommandHandler{
			ledger: ledger,
		})
		assert.Nil(t, err)

		_, err = cmdBus.Execute(&unsupportedCommand{})
		assert.IsType(t, &cq.UnsupportedRequestError{}, err)
	}

	// entity ids
	var accountId uuid.UUID
	var withdrawalId uuid.UUID

	// prepare command bus
	cmdBus := cq.CommandBus{}.New()
	err := cmdBus.RegisterHandler("command.create_account", &CreateAccountCommandHandler{
		ledger: ledger,
	})
	assert.Nil(t, err)

	err = cmdBus.RegisterHandler("command.create_planned_cash_withdrawal", &CreatePlannedCashWithdrawalCommandHandler{
		ledger: ledger,
	})
	assert.Nil(t, err)

	err = cmdBus.RegisterHandler("command.confirm_planned_cash_withdrawal", &ConfirmPlannedCashWithdrawalCommandHandler{
		ledger: ledger,
	})
	assert.Nil(t, err)

	// negative test for account not found error
	{
		confirmCommand := command.ConfirmPlannedCashWithdrawalCommand{}.New(uuid.NewV4(), uuid.NewV4())
		_, err := cmdBus.Execute(confirmCommand)
		assert.IsType(t, &accounting.AccountNotFoundError{}, err)
	}

	// create test account
	{
		createAccountCommand := command.CreateAccountCommand{}.New("Test Account", "EUR")
		id, err := cmdBus.Execute(createAccountCommand)
		assert.Nil(t, err)
		assert.IsType(t, uuid.UUID{}, id)

		err = ledger.AddValue(id.(uuid.UUID), accounting.Money{}.NewFromInt(100000, "EUR"), "initial")
		assert.Nil(t, err)

		accountId = id.(uuid.UUID)
	}

	// create planned cash withdrawal for account
	{
		createPlannedCashWithdrawalCommand := command.CreatePlannedCashWithdrawalCommand{}.New(
			accountId,
			(time.Now()).Add(time.Duration(time.Second*5)),
			accounting.Money{}.NewFromInt(10000, "EUR"), // 100.00 EUR
			"For testing only",
		)

		id, err := cmdBus.Execute(createPlannedCashWithdrawalCommand)
		assert.Nil(t, err)
		assert.IsType(t, uuid.UUID{}, id)

		withdrawalId = id.(uuid.UUID)
	}

	// confirm cash withdrawal
	{
		confirmCommand := command.ConfirmPlannedCashWithdrawalCommand{}.New(accountId, withdrawalId)
		_, err := cmdBus.Execute(confirmCommand)
		assert.Nil(t, err)
	}

	// load account and check balance
	{
		account, err := ledger.LoadAccount(accountId)
		assert.Nil(t, err)
		assert.NotNil(t, account)

		expectedBalance := accounting.Money{}.NewFromInt(90000, "EUR")
		assert.True(t, account.GetBalance().IsEqual(expectedBalance))
	}
}
