package accounting

import (
	"github.com/satori/go.uuid"
	"time"
)

// PlannedCashWithdrawalCreatedEvent
//
//
type PlannedCashWithdrawalCreatedEvent struct {
	AccountId uuid.UUID
	Date      time.Time
	Amount    Money
	Title     string
}

func (PlannedCashWithdrawalCreatedEvent) New(
	accountId uuid.UUID,
	date time.Time,
	amount Money,
	title string) *PlannedCashWithdrawalCreatedEvent {
	return &PlannedCashWithdrawalCreatedEvent{accountId, date, amount, title}
}

func (e *PlannedCashWithdrawalCreatedEvent) GetName() string {
	return "event.planned_cash_withdrawal_created"
}

func (e *PlannedCashWithdrawalCreatedEvent) GetAccountId() uuid.UUID {
	return e.AccountId
}
