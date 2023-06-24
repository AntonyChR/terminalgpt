package customErrors

import "fmt"

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

type InputError struct {
	Err   error
	Input string
}

func (i *InputError) Error() string {
	return fmt.Sprintf("Input %s: err %v", i.Input, i.Err)
}
