package errors

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// AccountNotFoundError
//
//
type AccountNotFoundError struct {
	AccountId uuid.UUID
}

func (e *AccountNotFoundError) Error() string {
	return fmt.Sprintf("Account with Id \"%s\" was not found.", e.AccountId.String())
}
