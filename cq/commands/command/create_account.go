package command

// CreateAccountCommand create a new account
//
//
type CreateAccountCommand struct {
	Title      string
	CurrencyId string
}

func (CreateAccountCommand) New(title, currencyId string) *CreateAccountCommand {
	cmd := &CreateAccountCommand{
		Title:      title,
		CurrencyId: currencyId,
	}

	return cmd
}

func (c *CreateAccountCommand) GetRequestId() string {
	return "command.create_account"
}
