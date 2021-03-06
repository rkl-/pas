package handler

import (
	"pas/src/accounting"
	"pas/src/cq"
	commandPkg "pas/src/cq/commands/command"
)

// ConfirmPlannedCashReceiptCommandHandler handler for ConfirmPlannedCashReceiptCommand
//
//
type ConfirmPlannedCashReceiptCommandHandler struct {
	ledger accounting.Ledger
}

func (h *ConfirmPlannedCashReceiptCommandHandler) Handle(request cq.Request) (interface{}, error) {
	command, ok := request.(*commandPkg.ConfirmPlannedCashReceiptCommand)
	if !ok {
		return nil, &cq.UnsupportedRequestError{}
	}

	if err := h.ledger.ConfirmPlannedCashReceipt(command.AccountId, command.ReceiptId); err != nil {
		return nil, err
	}

	return nil, nil
}
