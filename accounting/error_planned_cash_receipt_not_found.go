package accounting

import (
	"fmt"
	"github.com/satori/go.uuid"
)

type PlannedCashReceiptNotFoundError struct {
	receiptId uuid.UUID
	accountId uuid.UUID
}

func (e *PlannedCashReceiptNotFoundError) Error() string {
	return fmt.Sprintf("Planned cash receipt \"%s\" was not found for account \"%s\"", e.receiptId.String(), e.accountId.String())
}
