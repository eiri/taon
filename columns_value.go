package main

import "strings"

// ColumnsValue allows to parse comma separated list of columns
type ColumnsValue []string

// Set is a setter for ColumnsValue
func (c *ColumnsValue) Set(value string) error {
	parts := strings.Split(value, ",")
	*c = append(*c, parts...)
	return nil
}

func (c *ColumnsValue) String() string {
	return ""
}
