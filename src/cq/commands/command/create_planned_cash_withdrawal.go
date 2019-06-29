package command

import (
	"github.com/satori/go.uuid"
	"pas/src/money"
	"time"
)

// CreatePlannedCashWithdrawalCommand command which creates an income
//
//
type CreatePlannedCashWithdrawalCommand struct {
	BookingAccountId uuid.UUID
	Date             time.Time
	Amount           money.Money
	Title            string
}

func (c CreatePlannedCashWithdrawalCommand) New(
	bookingAccountId uuid.UUID,
	date time.Time,
	amount money.Money,
	title string) *CreatePlannedCashWithdrawalCommand {
	cmd := &CreatePlannedCashWithdrawalCommand{
		BookingAccountId: bookingAccountId,
		Date:             date,
		Amount:           amount,
		Title:            title,
	}

	return cmd
}

func (c *CreatePlannedCashWithdrawalCommand) GetRequestId() string {
	return "command.create_planned_cash_withdrawal"
}
