package cq

import "fmt"

// UnsupportedRequestError
//
//
type UnsupportedRequestError struct {
}

func (e *UnsupportedRequestError) Error() string {
	return fmt.Sprintf("unsupported request")
}
