package cq_command

import (
	"testing"
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
	//// prepare test account
	//ledger := accounting.Ledger{}.New(events.DomainDispatcher{}.GetInstance(), &events.InMemoryEventStorage{})
	//postBankAccount := ledger.CreateAccount("Postbank", "EUR")
	//
	//incomeAmount := accounting.Money{}.NewFromInt(1000000, "EUR") // 10,000.00 EUR
	//incomeTitle := "Salary"
	//
	//// prepare command bus
	//cmdBus := CommandBus{}.GetInstance()
	//err := cmdBus.RegisterHandler("command.create_planned_cash_receipt", &CreatePlannedCashReceiptCommandHandler{})
	//assert.Nil(t, err)
	//
	//// negative test for UnsupportedRequestError
	//_, err = cmdBus.Execute(&unsupportedCommand{})
	//assert.IsType(t, &UnsupportedRequestError{}, err)
	//
	//// negative test for DateInPastError
	//pastDate := time.Now()
	//pastDate.Add(time.Duration(-1) * time.Second)
	//
	//command := CreatePlannedCashReceiptCommand{}.New(postBankAccount.GetId(), pastDate, incomeAmount, incomeTitle)
	//_, err = cmdBus.Execute(command)
}
