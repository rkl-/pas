package errors

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// PlannedCashReceiptNotFoundError
//
//
type PlannedCashReceiptNotFoundError struct {
	ReceiptId uuid.UUID
	AccountId uuid.UUID
}

func (e *PlannedCashReceiptNotFoundError) Error() string {
	return fmt.Sprintf("Planned cash receipt \"%s\" was not found for account \"%s\"", e.ReceiptId.String(), e.AccountId.String())
}
