package cq_command

import (
	"pas/accounting"
	"pas/cq"
	"time"
)

// CreatePlannedCashWithdrawalCommandHandler handler for CreatePlannedCashWithdrawalCommand
//
//
type CreatePlannedCashWithdrawalCommandHandler struct {
	ledger accounting.Ledger
}

func (h *CreatePlannedCashWithdrawalCommandHandler) Handle(request cq.Request) (interface{}, error) {
	command, ok := request.(*CreatePlannedCashWithdrawalCommand)
	if !ok {
		return nil, &cq.UnsupportedRequestError{}
	}

	if time.Now().After(command.date) {
		return nil, &DateInPastError{}
	}

	if !h.ledger.HasAccount(command.bookingAccountId) {
		return nil, &AccountNotFoundError{command.bookingAccountId}
	}

	cashWithdrawal := accounting.PlannedCashFlow{}.New(command.date, command.amount, command.title)

	acc, err := h.ledger.LoadAccount(command.bookingAccountId)
	if err != nil {
		return nil, err
	}

	if err := h.ledger.AddPlannedCashWithdrawal(acc, cashWithdrawal); err != nil {
		return nil, err
	}

	return nil, nil
}
