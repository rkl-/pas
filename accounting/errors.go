package accounting

// InsufficientFoundsError
//
//
type InsufficientFoundsError struct {
}

func (e *InsufficientFoundsError) Error() string {
	return "insufficient founds"
}

// UnequalCurrenciesError
//
//
type UnequalCurrenciesError struct {
}

func (e *UnequalCurrenciesError) Error() string {
	return "unequal currencies"
}

// AccountCreatedEventNotFoundError
//
//
type AccountCreatedEventNotFoundError struct {
}

func (e *AccountCreatedEventNotFoundError) Error() string {
	return "initial account created event was not found"
}
