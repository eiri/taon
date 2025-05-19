package taon

import (
	"testing"
)

var cell_tests = []struct {
	cell   any
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
