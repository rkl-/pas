package command

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"pas/src/money"
	"testing"
	"time"
)

// TestCreatePlannedCashReceiptCommand_New
//
//
func TestCreatePlannedCashReceiptCommand_New(t *testing.T) {
	accountId := uuid.NewV4()
	date := time.Now()
	value := money.Money{}.NewFromInt(10000, "EUR") // 100.00 EUR
	title := "FooBar Title"

	pr := CreatePlannedCashReceiptCommand{}.New(accountId, date, value, title)
	assert.Equal(t, accountId, pr.BookingAccountId)
	assert.Equal(t, date, pr.Date)
	assert.Equal(t, value, pr.Amount)
	assert.Equal(t, title, pr.Title)
}

// TestCreatePlannedCashReceiptCommand_GetRequestId
//
//
func TestCreatePlannedCashReceiptCommand_GetRequestId(t *testing.T) {
	pr := &CreatePlannedCashReceiptCommand{}
	assert.Equal(t, "command.create_planned_cash_receipt", pr.GetRequestId())
}
