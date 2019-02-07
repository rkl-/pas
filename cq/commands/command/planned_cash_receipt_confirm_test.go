package command

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConfirmPlannedCashReceiptCommand_New
//
//
func TestConfirmPlannedCashReceiptCommand_New(t *testing.T) {
	accountId := uuid.NewV4()
	receiptId := uuid.NewV4()

	command := ConfirmPlannedCashReceiptCommand{}.New(accountId, receiptId)
	assert.IsType(t, &ConfirmPlannedCashReceiptCommand{}, command)
	assert.Equal(t, accountId, command.AccountId)
	assert.Equal(t, receiptId, command.ReceiptId)
}

// TestConfirmPlannedCashReceiptCommand_GetRequestId
//
//
func TestConfirmPlannedCashReceiptCommand_GetRequestId(t *testing.T) {
	command := &ConfirmPlannedCashReceiptCommand{}
	assert.Equal(t, "command.confirm_planned_cash_receipt", command.GetRequestId())
}
