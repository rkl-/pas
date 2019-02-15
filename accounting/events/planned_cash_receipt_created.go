package events

import (
	"github.com/satori/go.uuid"
	"pas/accounting/structs"
	"pas/money"
	"time"
)

// PlannedCashReceiptCreatedEvent
//
//
type PlannedCashReceiptCreatedEvent struct {
	ReceiptId uuid.UUID
	AccountId uuid.UUID
	Date      time.Time
	Amount    money.Money
	Title     string
}

func (PlannedCashReceiptCreatedEvent) NewFrom(flow *structs.PlannedCashFlow) *PlannedCashReceiptCreatedEvent {
	return &PlannedCashReceiptCreatedEvent{
		ReceiptId: flow.GetId(),
		AccountId: flow.AccountId,
		Date:      flow.Date,
		Amount:    flow.Amount,
		Title:     flow.Title,
	}
}

func (e *PlannedCashReceiptCreatedEvent) GetName() string {
	return "event.planned_cash_receipt_created"
}

func (e *PlannedCashReceiptCreatedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
