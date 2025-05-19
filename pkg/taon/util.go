package taon

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
)

// makeCell converts from typed input to string representation
func makeCell(v any) string {
	switch vv := v.(type) {
	case string:
		return vv
	case bool:
		return strconv.FormatBool(vv)
	case int:
		return strconv.Itoa(vv)
	case int64:
		return strconv.FormatInt(vv, 10)
	case uint64:
		return strconv.FormatUint(vv, 10)
	case float64:
		return strconv.FormatFloat(vv, 'f', 2, 64)
	case fmt.Stringer:
		return vv.String()
	}
	return fmt.Sprintf("%v", v)
}

// flatten recursively flattens nested maps and slices into dot-separated key paths
func flatten(topKey string, value any, result map[string]any) {
	switch v := value.(type) {
	case map[string]any:
		for key, val := range v {
			newKey := key
			if topKey != "" {
				newKey = topKey + "." + key
			}
			flatten(newKey, val, result)
		}
	case []any:
		for i, item := range v {
			newKey := topKey + "." + strconv.Itoa(i)
			flatten(newKey, item, result)
		}
	default:
		result[topKey] = v
	}
}

// AllocateColumnWidths distributes widths as evenly as possible across columns,
// trying to fully utilize maxWidth, without exceeding each column's desired width.
func AllocateColumnWidths(desiredWidths []int) []int {
	maxTableWidth, err := MaxTermWidth()
	if err != nil {
		maxTableWidth = 60
	}

	n := len(desiredWidths)
	margins := 3*n + 1
	maxWidth := maxTableWidth - margins

	allocated := make([]int, n)
	remainingWidths := make([]int, n)
	copy(remainingWidths, desiredWidths)

	remainingCols := make([]int, n)
	for i := range remainingCols {
		remainingCols[i] = i
	}

	remainingSpace := maxWidth

	for len(remainingCols) > 0 && remainingSpace > 0 {
		evenShare := remainingSpace / len(remainingCols)
		newRemainingCols := []int{}

		for _, i := range remainingCols {
			grant := min(evenShare, remainingWidths[i])
			allocated[i] += grant
			remainingWidths[i] -= grant
			remainingSpace -= grant

			if remainingWidths[i] > 0 {
				newRemainingCols = append(newRemainingCols, i)
			}
		}

		// No further distribution possible
		if len(newRemainingCols) == len(remainingCols) {
			break
		}

		remainingCols = newRemainingCols
	}

	// If there's still space left, distribute 1 unit at a time round-robin
	for i := 0; remainingSpace > 0; i = (i + 1) % n {
		if allocated[i] < desiredWidths[i] {
			allocated[i]++
			remainingSpace--
		}
		// If no more columns can grow, stop
		if i == n-1 && remainingSpace > 0 {
			canGrow := false
			for j := range n {
				if allocated[j] < desiredWidths[j] {
					canGrow = true
					break
				}
			}
			if !canGrow {
				break
			}
		}
	}

	return allocated
}

// MaxTermWidth returns tty's width
func MaxTermWidth() (int, error) {
	// prefer env $COLUMNS, fail back on tty if not set
	if columns := os.Getenv("COLUMNS"); columns != "" {
		return strconv.Atoi(columns)
	}
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	return width, err
}
