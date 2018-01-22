package errors

import (
	"errors"
	"fmt"
	"strings"
)

var TypeError = errors.New("type error")

type CompositeError struct {
	errors []error
}

func (c *CompositeError) Error() string {
	cs := make([]string, len(c.errors))
	for i, e := range c.errors {
		cs[i] = e.Error()
	}
	return fmt.Sprintf("composite error {%q}", strings.Join(cs, ", "))
}
func (c *CompositeError) Append(err error) {
	c.errors = append(c.errors, err)
}
