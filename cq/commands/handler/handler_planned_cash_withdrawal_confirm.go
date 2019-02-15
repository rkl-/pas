package handler

import (
	"pas/accounting"
	"pas/cq"
	commandPkg "pas/cq/commands/command"
)

// ConfirmPlannedCashWithdrawalCommandHandler handler for ConfirmPlannedCashWithdrawalCommand
//
//
type ConfirmPlannedCashWithdrawalCommandHandler struct {
	ledger accounting.Ledger
}

func (h *ConfirmPlannedCashWithdrawalCommandHandler) Handle(request cq.Request) (interface{}, error) {
	command, ok := request.(*commandPkg.ConfirmPlannedCashWithdrawalCommand)
	if !ok {
		return nil, &cq.UnsupportedRequestError{}
	}

	if err := h.ledger.ConfirmPlannedCashWithdrawal(command.AccountId, command.WithdrawalId); err != nil {
		return nil, err
	}

	return nil, nil
}
