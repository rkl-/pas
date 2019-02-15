package events

import (
	"github.com/satori/go.uuid"
	"pas/accounting/structs"
	"pas/money"
	"time"
)

// PlannedCashWithdrawalCreatedEvent
//
//
type PlannedCashWithdrawalCreatedEvent struct {
	WithdrawalId uuid.UUID
	AccountId    uuid.UUID
	Date         time.Time
	Amount       money.Money
	Title        string
}

func (PlannedCashWithdrawalCreatedEvent) NewFrom(flow *structs.PlannedCashFlow) *PlannedCashWithdrawalCreatedEvent {
	return &PlannedCashWithdrawalCreatedEvent{
		WithdrawalId: flow.GetId(),
		AccountId:    flow.AccountId,
		Date:         flow.Date,
		Amount:       flow.Amount,
		Title:        flow.Title,
	}
}

func (e *PlannedCashWithdrawalCreatedEvent) GetName() string {
	return "event.planned_cash_withdrawal_created"
}

func (e *PlannedCashWithdrawalCreatedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
