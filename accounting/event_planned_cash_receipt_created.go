package accounting

import (
	"github.com/satori/go.uuid"
	"time"
)

// PlannedCashReceiptCreatedEvent
//
//
type PlannedCashReceiptCreatedEvent struct {
	AccountId uuid.UUID
	Date      time.Time
	Amount    Money
	Title     string
}

func (PlannedCashReceiptCreatedEvent) New(
	accountId uuid.UUID,
	date time.Time,
	amount Money,
	title string) *PlannedCashReceiptCreatedEvent {
	return &PlannedCashReceiptCreatedEvent{accountId, date, amount, title}
}

func (e *PlannedCashReceiptCreatedEvent) GetName() string {
	return "event.planned_cash_receipt_created"
}

func (e *PlannedCashReceiptCreatedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
