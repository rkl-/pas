package handler

import (
	"pas/src/accounting"
	"pas/src/accounting/structs"
	"pas/src/cq"
	commandPkg "pas/src/cq/commands/command"
	"time"
)

// CreatePlannedCashReceiptCommandHandler handler for CreatePlannedCashReceiptCommand
//
//
type CreatePlannedCashReceiptCommandHandler struct {
	ledger accounting.Ledger
}

func (h *CreatePlannedCashReceiptCommandHandler) Handle(request cq.Request) (interface{}, error) {
	command, ok := request.(*commandPkg.CreatePlannedCashReceiptCommand)
	if !ok {
		return nil, &cq.UnsupportedRequestError{}
	}

	if time.Now().After(command.Date) {
		return nil, &DateInPastError{}
	}

	cashReceipt := structs.PlannedCashFlow{}.New(command.BookingAccountId, command.Date, command.Amount, command.Title)

	if err := h.ledger.AddPlannedCashReceipt(command.BookingAccountId, cashReceipt); err != nil {
		return nil, err
	}

	return cashReceipt.GetId(), nil
}
