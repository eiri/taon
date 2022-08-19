package taon

import (
	"reflect"
	"testing"
)

// TestColumns to complaiance to flag's Value type
func TestColumns(t *testing.T) {
	var c Columns
	if c.String() != "[]" {
		t.Errorf("Expecting `[]` for zero Columns")
	}
	c.Set("zebra,alpha,comma")
	expect := Columns{"zebra", "alpha", "comma"}
	if !reflect.DeepEqual(expect, c) {
		t.Errorf("Expecting %#v, got %#v", expect, c)
	}
	if c.String() != "[zebra alpha comma]" {
		t.Errorf("Expecting `[zebra alpha comma]` for Columns")
	}
}
