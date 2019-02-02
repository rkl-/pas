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

// TestCreatePlannedCashReceiptCommandHandler_Handle
//
//
func TestCreatePlannedCashReceiptCommandHandler_Handle(t *testing.T) {
	// prepare test account
	dispatcher := events.DomainDispatcher{}.New()
	ledger := accounting.DefaultLedger{}.New(dispatcher, &events.InMemoryEventStorage{})
	postBankAccount, _ := ledger.CreateAccount("Postbank", "EUR")

	incomeAmount := accounting.Money{}.NewFromInt(1000000, "EUR") // 10,000.00 EUR
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

		command := CreatePlannedCashReceiptCommand{}.New(postBankAccount.GetId(), pastDate, incomeAmount, incomeTitle)
		_, err = cmdBus.Execute(command)
		assert.IsType(t, &DateInPastError{}, err)
	}

	// negative test for AccountNotFoundError
	{
		validDate := (time.Now()).Add(time.Duration(1) * time.Hour)

		command := CreatePlannedCashReceiptCommand{}.New(uuid.NewV4(), validDate, incomeAmount, incomeTitle)
		_, err = cmdBus.Execute(command)
		assert.IsType(t, &AccountNotFoundError{}, err)
	}

	// positive test
	{
		var catchedEvent *accounting.PlannedCashReceiptCreatedEvent

		handler := &testEventHandler{}
		handler.dynamicHandle = func(event events.Event) {
			catchedEvent = event.(*accounting.PlannedCashReceiptCreatedEvent)
		}

		dispatcher.RegisterHandler((&accounting.PlannedCashReceiptCreatedEvent{}).GetName(), handler)

		validDate := (time.Now()).Add(time.Duration(1) * time.Hour)

		command := CreatePlannedCashReceiptCommand{}.New(postBankAccount.GetId(), validDate, incomeAmount, incomeTitle)
		_, err = cmdBus.Execute(command)
		assert.Nil(t, err)
		assert.NotNil(t, catchedEvent)
		assert.Equal(t, catchedEvent.AccountId, postBankAccount.GetId())
		assert.Equal(t, catchedEvent.Title, incomeTitle)
		assert.Equal(t, catchedEvent.Amount, incomeAmount)
		assert.Equal(t, catchedEvent.Date, validDate)
	}
}
