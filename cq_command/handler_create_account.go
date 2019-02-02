package cq_command

import (
	"pas/accounting"
	"pas/cq"
)

// CreateAccountCommandHandler
//
//
type CreateAccountCommandHandler struct {
	ledger accounting.Ledger
}

func (c *CreateAccountCommandHandler) Handle(request cq.Request) (interface{}, error) {
	command, ok := request.(*CreateAccountCommand)
	if !ok {
		return nil, &cq.UnsupportedRequestError{}
	}

	account, err := c.ledger.CreateAccount(command.title, command.currencyId)
	if err != nil {
		return nil, err
	}

	return account.GetId(), nil
}
