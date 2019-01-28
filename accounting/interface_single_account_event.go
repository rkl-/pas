package accounting

import "github.com/satori/go.uuid"

// SingleAccountEvent event which has an unique account association
//
//
type SingleAccountEvent interface {
	GetAccountId() uuid.UUID
}
