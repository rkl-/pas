package cq_command

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/accounting"
	"pas/cq"
	"pas/events"
	"testing"
	"time"
)

// TestCreatePlannedCashWithdrawalCommandHandler_Handle
//
//
func TestCreatePlannedCashWithdrawalCommandHandler_Handle(t *testing.T) {
	// prepare test account
	dispatcher := events.DomainDispatcher{}.New()
	ledger := accounting.DefaultLedger{}.New(dispatcher, &events.InMemoryEventStorage{})
	postBankAccount, _ := ledger.CreateAccount("Postbank", "EUR")

	expenseAmount := accounting.Money{}.NewFromInt(1000000, "EUR") // 10,000.00 EUR
	expenseTitle := "Salary"

	// prepare command bus
	cmdBus := cq.CommandBus{}.New()
	err := cmdBus.RegisterHandler("command.unsupported_command", &CreatePlannedCashWithdrawalCommandHandler{
		ledger: ledger,
	})
	assert.Nil(t, err)

	err = cmdBus.RegisterHandler("command.create_planned_cash_withdrawal", &CreatePlannedCashWithdrawalCommandHandler{
		ledger: ledger,
	})
	assert.Nil(t, err)

	// negative test for UnsupportedRequestError
	{
		_, err = cmdBus.Execute(&unsupportedCommand{})
		assert.IsType(t, &cq.UnsupportedRequestError{}, err)
	}

	// negative test for DateInPastError
	{
		pastDate := (time.Now()).Add(time.Duration(-1) * time.Second)

		command := CreatePlannedCashWithdrawalCommand{}.New(postBankAccount.GetId(), pastDate, expenseAmount, expenseTitle)
		_, err = cmdBus.Execute(command)
		assert.IsType(t, &DateInPastError{}, err)
	}

	// negative test for AccountNotFoundError
	{
		validDate := (time.Now()).Add(time.Duration(1) * time.Hour)

		command := CreatePlannedCashWithdrawalCommand{}.New(uuid.NewV4(), validDate, expenseAmount, expenseTitle)
		_, err = cmdBus.Execute(command)
		assert.IsType(t, &AccountNotFoundError{}, err)
	}

	// positive test
	{
		var catchedEvent *accounting.PlannedCashWithdrawalCreatedEvent

		handler := &testEventHandler{}
		handler.dynamicHandle = func(event events.Event) {
			catchedEvent = event.(*accounting.PlannedCashWithdrawalCreatedEvent)
		}

		dispatcher.RegisterHandler((&accounting.PlannedCashWithdrawalCreatedEvent{}).GetName(), handler)

		validDate := (time.Now()).Add(time.Duration(1) * time.Hour)

		command := CreatePlannedCashWithdrawalCommand{}.New(postBankAccount.GetId(), validDate, expenseAmount, expenseTitle)
		_, err = cmdBus.Execute(command)
		assert.Nil(t, err)
		assert.NotNil(t, catchedEvent)
		assert.Equal(t, catchedEvent.AccountId, postBankAccount.GetId())
		assert.Equal(t, catchedEvent.Title, expenseTitle)
		assert.Equal(t, catchedEvent.Amount, expenseAmount)
		assert.Equal(t, catchedEvent.Date, validDate)
	}
}
