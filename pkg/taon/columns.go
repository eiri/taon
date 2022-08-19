package taon

import (
	"fmt"
	"strings"
)

// Columns allows to parse comma separated list of columns
type Columns []string

// Set is a setter for Columns
func (c *Columns) Set(value string) error {
	parts := strings.Split(value, ",")
	*c = append(*c, parts...)
	return nil
}

func (c *Columns) String() string {
	if c == nil {
		return "[]"
	}
	return fmt.Sprintf("[%s]", strings.Join(*c, " "))
}
