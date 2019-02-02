package cq_command

// CreateAccountCommand create a new account
//
//
type CreateAccountCommand struct {
	title      string
	currencyId string
}

func (CreateAccountCommand) New(title, currencyId string) *CreateAccountCommand {
	cmd := &CreateAccountCommand{
		title:      title,
		currencyId: currencyId,
	}

	return cmd
}

func (c *CreateAccountCommand) GetRequestId() string {
	return "command.create_account"
}
