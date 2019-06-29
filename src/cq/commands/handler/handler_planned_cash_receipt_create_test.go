package handler

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/src/accounting"
	"pas/src/accounting/errors"
	events2 "pas/src/accounting/events"
	"pas/src/cq"
	commandPkg "pas/src/cq/commands/command"
	"pas/src/events"
	"pas/src/money"
	"testing"
	"time"
)

// TestCreatePlannedCashReceiptCommandHandler_Handle
//
//
func TestCreatePlannedCashReceiptCommandHandler_Handle(t *testing.T) {
	// prepare test account
	dispatcher := events.DomainDispatcher{}.New()
	ledger := accounting.DefaultLedger{}.New(dispatcher, &events.InMemoryEventStorage{})
	postBankAccount, _ := ledger.CreateAccount("Postbank", "EUR")

	incomeAmount := money.Money{}.NewFromInt(1000000, "EUR") // 10,000.00 EUR
	incomeTitle := "Salary"

	// prepare command bus
	cmdBus := cq.CommandBus{}.New()
	err := cmdBus.RegisterHandler("command.unsupported_command", &CreatePlannedCashReceiptCommandHandler{
		ledger: ledger,
	})
	assert.Nil(t, err)

	err = cmdBus.RegisterHandler("command.create_planned_cash_receipt", &CreatePlannedCashReceiptCommandHandler{
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

		command := commandPkg.CreatePlannedCashReceiptCommand{}.New(postBankAccount.GetId(), pastDate, incomeAmount, incomeTitle)
		_, err = cmdBus.Execute(command)
		assert.IsType(t, &DateInPastError{}, err)
	}

	// negative test for AccountNotFoundError
	{
		validDate := (time.Now()).Add(time.Duration(1) * time.Hour)

		command := commandPkg.CreatePlannedCashReceiptCommand{}.New(uuid.NewV4(), validDate, incomeAmount, incomeTitle)
		_, err = cmdBus.Execute(command)
		assert.IsType(t, &errors.AccountNotFoundError{}, err)
	}

	// positive test
	{
		var catchedEvent *events2.PlannedCashReceiptCreatedEvent

		handler := &testEventHandler{}
		handler.dynamicHandle = func(event events.Event) {
			catchedEvent = event.(*events2.PlannedCashReceiptCreatedEvent)
		}

		dispatcher.RegisterHandler((&events2.PlannedCashReceiptCreatedEvent{}).GetName(), handler)

		validDate := (time.Now()).Add(time.Duration(1) * time.Hour)

		command := commandPkg.CreatePlannedCashReceiptCommand{}.New(postBankAccount.GetId(), validDate, incomeAmount, incomeTitle)
		id, err := cmdBus.Execute(command)
		assert.Nil(t, err)
		assert.IsType(t, uuid.UUID{}, id)
		assert.NotNil(t, catchedEvent)
		assert.Equal(t, catchedEvent.ReceiptId, id)
		assert.Equal(t, catchedEvent.AccountId, postBankAccount.GetId())
		assert.Equal(t, catchedEvent.Title, incomeTitle)
		assert.Equal(t, catchedEvent.Amount, incomeAmount)
		assert.Equal(t, catchedEvent.Date, validDate)
	}
}
