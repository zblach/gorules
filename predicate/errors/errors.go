package errors

import (
	"errors"
	"fmt"
	"strings"
)

var TypeError = errors.New("type error")

type CompositeError struct {
	Errors []error
}

func (c *CompositeError) Error() string {
	cs := make([]string, len(c.Errors))
	for i, e := range c.Errors {
		cs[i] = e.Error()
	}
	return fmt.Sprintf("composite error {%q}", strings.Join(cs, ", "))
}
func (c *CompositeError) Append(err ...error) {
	c.Errors = append(c.Errors, err...)
}

func (c *CompositeError) NilZero() error {
	if len(c.Errors) == 0 {
		return nil
	}
	return c
}
