package taon

import (
	"testing"
)

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
