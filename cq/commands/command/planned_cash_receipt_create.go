package command

import (
	"github.com/satori/go.uuid"
	"pas/money"
	"time"
)

// CreatePlannedCashReceiptCommand command which creates an income
//
//
type CreatePlannedCashReceiptCommand struct {
	BookingAccountId uuid.UUID
	Date             time.Time
	Amount           money.Money
	Title            string
}

func (c CreatePlannedCashReceiptCommand) New(
	bookingAccountId uuid.UUID,
	date time.Time,
	amount money.Money,
	title string) *CreatePlannedCashReceiptCommand {
	cmd := &CreatePlannedCashReceiptCommand{
		BookingAccountId: bookingAccountId,
		Date:             date,
		Amount:           amount,
		Title:            title,
	}

	return cmd
}

func (c *CreatePlannedCashReceiptCommand) GetRequestId() string {
	return "command.create_planned_cash_receipt"
}
