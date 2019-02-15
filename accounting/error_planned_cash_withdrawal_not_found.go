package accounting

import (
	"fmt"
	"github.com/satori/go.uuid"
)

type PlannedCashWithdrawalNotFoundError struct {
	withdrawalId uuid.UUID
	accountId    uuid.UUID
}

func (e *PlannedCashWithdrawalNotFoundError) Error() string {
	return fmt.Sprintf("Planned cash withdrawal \"%s\" was not found for account \"%s\"", e.withdrawalId.String(), e.accountId.String())
}
