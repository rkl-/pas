package accounting

// InsufficientFoundsError
//
//
type InsufficientFoundsError struct {
}

// Error implements error interface
//
//
func (e *InsufficientFoundsError) Error() string {
	return "insufficient founds"
}

// UnequalCurrenciesError
//
//
type UnequalCurrenciesError struct {
}

// Error implements error interface
//
//
func (e *UnequalCurrenciesError) Error() string {
	return "unequal currencies"
}
