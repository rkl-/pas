package command

import "github.com/satori/go.uuid"

// ConfirmPlannedCashWithdrawalCommand
//
//
type ConfirmPlannedCashWithdrawalCommand struct {
	AccountId    uuid.UUID
	WithdrawalId uuid.UUID
}

func (ConfirmPlannedCashWithdrawalCommand) New(accountId, withdrawalId uuid.UUID) *ConfirmPlannedCashWithdrawalCommand {
	return &ConfirmPlannedCashWithdrawalCommand{accountId, withdrawalId}
}

func (c *ConfirmPlannedCashWithdrawalCommand) GetRequestId() string {
	return "command.confirm_planned_cash_withdrawal"
}
