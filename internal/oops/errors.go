package oops

import (
	"fmt"
)

type CodedError struct {
	HumanMessage string
	Context      string
	HttpStatus   int
	Err          error
}

func (c CodedError) Error() string {
	m := fmt.Sprintf("%s (%d)", c.Context, c.HttpStatus)
	if c.Err != nil {
		return fmt.Sprintf("%s : %v", m, c.Err)
	}
	return m
}
