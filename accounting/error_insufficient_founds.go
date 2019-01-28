package accounting

// InsufficientFoundsError
//
//
type InsufficientFoundsError struct {
}

func (e *InsufficientFoundsError) Error() string {
	return "insufficient founds"
}
