package command

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConfirmPlannedCashWithdrawalCommand_New
//
//
func TestConfirmPlannedCashWithdrawalCommand_New(t *testing.T) {
	accountId := uuid.NewV4()
	withdrawalId := uuid.NewV4()

	command := ConfirmPlannedCashWithdrawalCommand{}.New(accountId, withdrawalId)
	assert.IsType(t, &ConfirmPlannedCashWithdrawalCommand{}, command)
	assert.Equal(t, accountId, command.AccountId)
	assert.Equal(t, withdrawalId, command.WithdrawalId)
}

// TestConfirmPlannedCashWithdrawalCommand_GetRequestId
//
//
func TestConfirmPlannedCashWithdrawalCommand_GetRequestId(t *testing.T) {
	command := &ConfirmPlannedCashWithdrawalCommand{}
	assert.Equal(t, "command.confirm_planned_cash_withdrawal", command.GetRequestId())
}
