package accounting

import "time"

// PlannedCashReceipt
//
//
type PlannedCashReceipt struct {
	date   time.Time
	amount Money
	title  string
}

func (PlannedCashReceipt) New(date time.Time, amount Money, title string) *PlannedCashReceipt {
	return &PlannedCashReceipt{date, amount, title}
}
