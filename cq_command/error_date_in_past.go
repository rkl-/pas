package cq_command

// DateInPastError
//
//
type DateInPastError struct {
}

func (e *DateInPastError) Error() string {
	return "date in past"
}
