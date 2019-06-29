package handler

import (
	"pas/src/accounting"
	"pas/src/accounting/structs"
	"pas/src/cq"
	commandPkg "pas/src/cq/commands/command"
	"time"
)

// CreatePlannedCashWithdrawalCommandHandler handler for CreatePlannedCashWithdrawalCommand
//
//
type CreatePlannedCashWithdrawalCommandHandler struct {
	ledger accounting.Ledger
}

func (h *CreatePlannedCashWithdrawalCommandHandler) Handle(request cq.Request) (interface{}, error) {
	command, ok := request.(*commandPkg.CreatePlannedCashWithdrawalCommand)
	if !ok {
		return nil, &cq.UnsupportedRequestError{}
	}

	if time.Now().After(command.Date) {
		return nil, &DateInPastError{}
	}

	cashWithdrawal := structs.PlannedCashFlow{}.New(command.BookingAccountId, command.Date, command.Amount, command.Title)

	if err := h.ledger.AddPlannedCashWithdrawal(command.BookingAccountId, cashWithdrawal); err != nil {
		return nil, err
	}

	return cashWithdrawal.GetId(), nil
}
