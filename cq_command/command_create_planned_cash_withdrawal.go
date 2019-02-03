package cq_command

import (
	"github.com/satori/go.uuid"
	"pas/accounting"
	"time"
)

// CreatePlannedCashWithdrawalCommand command which creates an income
//
//
type CreatePlannedCashWithdrawalCommand struct {
	bookingAccountId uuid.UUID
	date             time.Time
	amount           accounting.Money
	title            string
}

func (c CreatePlannedCashWithdrawalCommand) New(
	bookingAccountId uuid.UUID,
	date time.Time,
	amount accounting.Money,
	title string) *CreatePlannedCashWithdrawalCommand {
	cmd := &CreatePlannedCashWithdrawalCommand{
		bookingAccountId: bookingAccountId,
		date:             date,
		amount:           amount,
		title:            title,
	}

	return cmd
}

func (c *CreatePlannedCashWithdrawalCommand) GetRequestId() string {
	return "command.create_planned_cash_withdrawal"
}
