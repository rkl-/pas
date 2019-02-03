package accounting

import "time"

// PlannedCashFlow
//
//
type PlannedCashFlow struct {
	date   time.Time
	amount Money
	title  string
}

func (PlannedCashFlow) New(date time.Time, amount Money, title string) *PlannedCashFlow {
	return &PlannedCashFlow{date, amount, title}
}
