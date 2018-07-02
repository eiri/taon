package main

import (
	"reflect"
	"testing"
)

// TestColumnsValue to complaiance to flag's Value type
func TestColumnsValue(t *testing.T) {
	var c ColumnsValue
	if c.String() != "[]" {
		t.Errorf("Expecting `[]` for zero ColumnsValue")
	}
	c.Set("zebra,alpha,comma")
	expect := ColumnsValue{"zebra", "alpha", "comma"}
	if !reflect.DeepEqual(expect, c) {
		t.Errorf("Expecting %#v, got %#v", expect, c)
	}
	if c.String() != "[zebra alpha comma]" {
		t.Errorf("Expecting `[zebra alpha comma]` for ColumnsValue")
	}
}
