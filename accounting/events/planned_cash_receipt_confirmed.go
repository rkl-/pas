package events

import (
	"github.com/satori/go.uuid"
	"pas/accounting/structs"
	"pas/money"
	"time"
)

// PlannedCashReceiptConfirmedEvent
//
//
type PlannedCashReceiptConfirmedEvent struct {
	ReceiptId uuid.UUID
	AccountId uuid.UUID
	Date      time.Time
	Amount    money.Money
	Title     string
}

func (PlannedCashReceiptConfirmedEvent) NewFrom(flow *structs.PlannedCashFlow) *PlannedCashReceiptConfirmedEvent {
	return &PlannedCashReceiptConfirmedEvent{
		ReceiptId: flow.GetId(),
		AccountId: flow.AccountId,
		Date:      flow.Date,
		Amount:    flow.Amount,
		Title:     flow.Title,
	}
}

func (e *PlannedCashReceiptConfirmedEvent) GetName() string {
	return "event.planned_cash_receipt_confirmed"
}

func (e *PlannedCashReceiptConfirmedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
