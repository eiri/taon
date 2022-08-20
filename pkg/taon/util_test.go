package taon

import (
	"reflect"
	"testing"
)

var header_tests = []struct {
	columns Columns
	expect  Header
}{
	{Columns{}, Header{"a", "b", "o", "r", "z"}},
	{Columns{"r", "a"}, Header{"r", "a"}},
	{Columns{"b", "x", "z", "p"}, Header{"b", "z"}},
}

// TestMakeHeader to ensure we are getting sorted list of strings
func TestMakeHeader(t *testing.T) {
	in := map[string]int{"b": 1, "z": 2, "a": 3, "r": 4, "o": 5}
	for _, tt := range header_tests {
		header, err := makeHeader(in, tt.columns)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tt.expect, header) {
			t.Errorf("Expecting %#v, got %#v", tt.expect, header)
		}
	}

	// return error
	_, err := makeHeader(in, Columns{"none"})
	if err == nil {
		t.Error("Expecting error, got nil")
	}
}

var cell_tests = []struct {
	cell   interface{}
	expect string
}{
	{"a", "a"},
	{"long string with ⭐️", "long string with ⭐️"},
	{'a', "97"},
	{1, "1"},
	{int64(1), "1"},
	{uint64(1), "1"},
	{1.1, "1.10"},
	{1.12, "1.12"},
	{1.123, "1.12"},
	{true, "true"},
	{false, "false"},
}

// TestMakeHeader to ensure we are getting sorted list of strings
func TestMakeCell(t *testing.T) {
	for _, tt := range cell_tests {
		cell := makeCell(tt.cell)
		if cell != tt.expect {
			t.Errorf("Expecting %q, got %q", tt.expect, cell)
		}
	}
}

var max_tests = []struct {
	a, b, c, expect int
}{
	{1, 2, 3, 3},
	{2, 1, 3, 3},
	{3, 2, 1, 3},
	{0, 0, 1, 1},
	{5, 5, 5, 5},
}

// TestMakeHeader to ensure we are getting sorted list of strings
func TestMax(t *testing.T) {
	for _, tt := range max_tests {
		m := max(tt.a, tt.b, tt.c)
		if m != tt.expect {
			t.Errorf("Expecting %d, got %d", tt.expect, m)
		}
	}
}
