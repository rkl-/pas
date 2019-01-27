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

type unsupportedCommand struct {
}

func (c *unsupportedCommand) GetRequestId() string {
	return "command.create_planned_cash_receipt"
}

// TestCreatePlannedCashReceiptCommandHandler_Handle
//
//
func TestCreatePlannedCashReceiptCommandHandler_Handle(t *testing.T) {
	// prepare test account
	ledger := accounting.DefaultLedger{}.New(events.DomainDispatcher{}.New(), &events.InMemoryEventStorage{})
	postBankAccount := ledger.CreateAccount("Postbank", "EUR")

	incomeAmount := accounting.Money{}.NewFromInt(1000000, "EUR") // 10,000.00 EUR
	incomeTitle := "Salary"

	// prepare dispatcher
	dispatcher := events.DomainDispatcher{}.New()

	// prepare command bus
	cmdBus := cq.CommandBus{}.New()
	err := cmdBus.RegisterHandler("command.create_planned_cash_receipt", &CreatePlannedCashReceiptCommandHandler{
		dispatcher: dispatcher,
		ledger:     ledger,
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
}
