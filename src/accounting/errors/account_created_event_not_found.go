package errors

// AccountCreatedEventNotFoundError
//
//
type AccountCreatedEventNotFoundError struct {
}

func (e *AccountCreatedEventNotFoundError) Error() string {
	return "initial account created event was not found"
}
