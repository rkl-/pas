package accounting

import (
	"github.com/satori/go.uuid"
	"time"
)

type PlannedCashFlowMap map[uuid.UUID]*PlannedCashFlow

// PlannedCashFlow
//
//
type PlannedCashFlow struct {
	id        uuid.UUID
	accountId uuid.UUID
	date      time.Time
	amount    Money
	title     string
}

func (PlannedCashFlow) New(accountId uuid.UUID, date time.Time, amount Money, title string) *PlannedCashFlow {
	return &PlannedCashFlow{uuid.NewV4(), accountId, date, amount, title}
}

func (f *PlannedCashFlow) GetId() uuid.UUID {
	return f.id
}
