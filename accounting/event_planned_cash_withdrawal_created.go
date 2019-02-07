package accounting

import (
	"github.com/satori/go.uuid"
	"time"
)

// PlannedCashWithdrawalCreatedEvent
//
//
type PlannedCashWithdrawalCreatedEvent struct {
	WithdrawalId uuid.UUID
	AccountId    uuid.UUID
	Date         time.Time
	Amount       Money
	Title        string
}

func (PlannedCashWithdrawalCreatedEvent) NewFrom(flow *PlannedCashFlow) *PlannedCashWithdrawalCreatedEvent {
	return &PlannedCashWithdrawalCreatedEvent{
		WithdrawalId: flow.GetId(),
		AccountId:    flow.accountId,
		Date:         flow.date,
		Amount:       flow.amount,
		Title:        flow.title,
	}
}

func (e *PlannedCashWithdrawalCreatedEvent) GetName() string {
	return "event.planned_cash_withdrawal_created"
}

func (e *PlannedCashWithdrawalCreatedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
