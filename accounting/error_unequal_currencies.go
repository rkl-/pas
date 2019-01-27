package accounting

// UnequalCurrenciesError
//
//
type UnequalCurrenciesError struct {
}

func (e *UnequalCurrenciesError) Error() string {
	return "unequal currencies"
}
