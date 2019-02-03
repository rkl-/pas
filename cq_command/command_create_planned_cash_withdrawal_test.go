package cq_command

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/accounting"
	"testing"
	"time"
)

// TestCreatePlannedCashWithdrawalCommand_New
//
//
func TestCreatePlannedCashWithdrawalCommand_New(t *testing.T) {
	accountId := uuid.NewV4()
	date := time.Now()
	value := accounting.Money{}.NewFromInt(10000, "EUR") // 100.00 EUR
	title := "FooBar Title"

	pr := CreatePlannedCashWithdrawalCommand{}.New(accountId, date, value, title)
	assert.Equal(t, accountId, pr.bookingAccountId)
	assert.Equal(t, date, pr.date)
	assert.Equal(t, value, pr.amount)
	assert.Equal(t, title, pr.title)
}

// TestCreatePlannedCashWithdrawalCommand_GetRequestId
//
//
func TestCreatePlannedCashWithdrawalCommand_GetRequestId(t *testing.T) {
	pr := &CreatePlannedCashWithdrawalCommand{}
	assert.Equal(t, "command.create_planned_cash_withdrawal", pr.GetRequestId())
}
