package accounting

import (
	"github.com/satori/go.uuid"
	"time"
)

// PlannedCashReceiptCreatedEvent
//
//
type PlannedCashReceiptCreatedEvent struct {
	ReceiptId uuid.UUID
	AccountId uuid.UUID
	Date      time.Time
	Amount    Money
	Title     string
}

func (PlannedCashReceiptCreatedEvent) NewFrom(flow *PlannedCashFlow) *PlannedCashReceiptCreatedEvent {
	return &PlannedCashReceiptCreatedEvent{
		ReceiptId: flow.GetId(),
		AccountId: flow.accountId,
		Date:      flow.date,
		Amount:    flow.amount,
		Title:     flow.title,
	}
}

func (e *PlannedCashReceiptCreatedEvent) GetName() string {
	return "event.planned_cash_receipt_created"
}

func (e *PlannedCashReceiptCreatedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
