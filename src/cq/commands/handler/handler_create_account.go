package handler

import (
	"pas/src/accounting"
	"pas/src/cq"
	commandPkg "pas/src/cq/commands/command"
)

// CreateAccountCommandHandler
//
//
type CreateAccountCommandHandler struct {
	ledger accounting.Ledger
}

func (c *CreateAccountCommandHandler) Handle(request cq.Request) (interface{}, error) {
	command, ok := request.(*commandPkg.CreateAccountCommand)
	if !ok {
		return nil, &cq.UnsupportedRequestError{}
	}

	account, err := c.ledger.CreateAccount(command.Title, command.CurrencyId)
	if err != nil {
		return nil, err
	}

	return account.GetId(), nil
}
