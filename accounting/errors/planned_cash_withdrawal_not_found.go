package errors

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// PlannedCashWithdrawalNotFoundError
//
//
type PlannedCashWithdrawalNotFoundError struct {
	WithdrawalId uuid.UUID
	AccountId    uuid.UUID
}

func (e *PlannedCashWithdrawalNotFoundError) Error() string {
	return fmt.Sprintf("Planned cash withdrawal \"%s\" was not found for account \"%s\"", e.WithdrawalId.String(), e.AccountId.String())
}
