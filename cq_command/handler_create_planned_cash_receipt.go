package cq_command

import (
	"pas/accounting"
	"pas/cq"
	"time"
)

// CreatePlannedCashReceiptCommandHandler handler for CreatePlannedCashReceiptCommand
//
//
type CreatePlannedCashReceiptCommandHandler struct {
	ledger accounting.Ledger
}

func (h *CreatePlannedCashReceiptCommandHandler) Handle(request cq.Request) (interface{}, error) {
	command, ok := request.(*CreatePlannedCashReceiptCommand)
	if !ok {
		return nil, &cq.UnsupportedRequestError{}
	}

	if time.Now().After(command.date) {
		return nil, &DateInPastError{}
	}

	if !h.ledger.HasAccount(command.bookingAccountId) {
		return nil, &AccountNotFoundError{command.bookingAccountId}
	}

	cashReceipt := accounting.PlannedCashReceipt{}.New(command.date, command.amount, command.title)

	acc, err := h.ledger.LoadAccount(command.bookingAccountId)
	if err != nil {
		return nil, err
	}

	if err := h.ledger.AddPlannedCashReceipt(acc, cashReceipt); err != nil {
		return nil, err
	}

	return nil, nil
}
