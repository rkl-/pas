package accounting

import (
	"github.com/satori/go.uuid"
	"time"
)

// PlannedCashWithdrawalConfirmedEvent
//
//
type PlannedCashWithdrawalConfirmedEvent struct {
	WithdrawalId uuid.UUID
	AccountId    uuid.UUID
	Date         time.Time
	Amount       Money
	Title        string
}

func (PlannedCashWithdrawalConfirmedEvent) NewFrom(flow *PlannedCashFlow) *PlannedCashWithdrawalConfirmedEvent {
	return &PlannedCashWithdrawalConfirmedEvent{
		WithdrawalId: flow.GetId(),
		AccountId:    flow.accountId,
		Date:         flow.date,
		Amount:       flow.amount,
		Title:        flow.title,
	}
}

func (e *PlannedCashWithdrawalConfirmedEvent) GetName() string {
	return "event.planned_cash_withdrawal_confirmed"
}

func (e *PlannedCashWithdrawalConfirmedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
