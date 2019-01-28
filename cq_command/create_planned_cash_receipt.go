package cq_command

import (
	"github.com/satori/go.uuid"
	"pas/accounting"
	"pas/cq"
	"time"
)

// CreatePlannedCashReceiptCommand command which creates an income
//
//
type CreatePlannedCashReceiptCommand struct {
	bookingAccountId uuid.UUID
	date             time.Time
	amount           accounting.Money
	title            string
}

func (c CreatePlannedCashReceiptCommand) New(
	bookingAccountId uuid.UUID,
	date time.Time,
	amount accounting.Money,
	title string) *CreatePlannedCashReceiptCommand {
	cmd := &CreatePlannedCashReceiptCommand{
		bookingAccountId: bookingAccountId,
		date:             date,
		amount:           amount,
		title:            title,
	}

	return cmd
}

func (c *CreatePlannedCashReceiptCommand) GetRequestId() string {
	return "command.create_planned_cash_receipt"
}

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
