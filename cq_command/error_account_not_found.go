package cq_command

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// AccountNotFoundError
//
//
type AccountNotFoundError struct {
	accountId uuid.UUID
}

func (e *AccountNotFoundError) Error() string {
	return fmt.Sprintf("Account with id \"%s\" was not found.", e.accountId.String())
}
