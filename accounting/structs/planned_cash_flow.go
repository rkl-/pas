package structs

import (
	"github.com/satori/go.uuid"
	"pas/money"
	"time"
)

type PlannedCashFlowMap map[uuid.UUID]*PlannedCashFlow

// PlannedCashFlow
//
//
type PlannedCashFlow struct {
	Id        uuid.UUID
	AccountId uuid.UUID
	Date      time.Time
	Amount    money.Money
	Title     string
}

func (PlannedCashFlow) New(accountId uuid.UUID, date time.Time, amount money.Money, title string) *PlannedCashFlow {
	return &PlannedCashFlow{uuid.NewV4(), accountId, date, amount, title}
}

func (f *PlannedCashFlow) GetId() uuid.UUID {
	return f.Id
}
