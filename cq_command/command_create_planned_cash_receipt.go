package cq_command

import (
	"github.com/satori/go.uuid"
	"pas/accounting"
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
