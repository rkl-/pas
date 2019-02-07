package command

import "github.com/satori/go.uuid"

// ConfirmPlannedCashReceiptCommand
//
//
type ConfirmPlannedCashReceiptCommand struct {
	AccountId uuid.UUID
	ReceiptId uuid.UUID
}

func (ConfirmPlannedCashReceiptCommand) New(accountId, receiptId uuid.UUID) *ConfirmPlannedCashReceiptCommand {
	return &ConfirmPlannedCashReceiptCommand{accountId, receiptId}
}

func (c *ConfirmPlannedCashReceiptCommand) GetRequestId() string {
	return "command.confirm_planned_cash_receipt"
}
