package events

import (
	"github.com/satori/go.uuid"
	"pas/accounting/structs"
	"pas/money"
	"time"
)

// PlannedCashWithdrawalConfirmedEvent
//
//
type PlannedCashWithdrawalConfirmedEvent struct {
	WithdrawalId uuid.UUID
	AccountId    uuid.UUID
	Date         time.Time
	Amount       money.Money
	Title        string
}

func (PlannedCashWithdrawalConfirmedEvent) NewFrom(flow *structs.PlannedCashFlow) *PlannedCashWithdrawalConfirmedEvent {
	return &PlannedCashWithdrawalConfirmedEvent{
		WithdrawalId: flow.GetId(),
		AccountId:    flow.AccountId,
		Date:         flow.Date,
		Amount:       flow.Amount,
		Title:        flow.Title,
	}
}

func (e *PlannedCashWithdrawalConfirmedEvent) GetName() string {
	return "event.planned_cash_withdrawal_confirmed"
}

func (e *PlannedCashWithdrawalConfirmedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
