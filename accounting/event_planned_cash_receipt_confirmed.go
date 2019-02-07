package accounting

import (
	"github.com/satori/go.uuid"
	"time"
)

// PlannedCashReceiptConfirmedEvent
//
//
type PlannedCashReceiptConfirmedEvent struct {
	ReceiptId uuid.UUID
	AccountId uuid.UUID
	Date      time.Time
	Amount    Money
	Title     string
}

func (PlannedCashReceiptConfirmedEvent) NewFrom(flow *PlannedCashFlow) *PlannedCashReceiptConfirmedEvent {
	return &PlannedCashReceiptConfirmedEvent{
		ReceiptId: flow.GetId(),
		AccountId: flow.accountId,
		Date:      flow.date,
		Amount:    flow.amount,
		Title:     flow.title,
	}
}

func (e *PlannedCashReceiptConfirmedEvent) GetName() string {
	return "event.planned_cash_receipt_confirmed"
}

func (e *PlannedCashReceiptConfirmedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
