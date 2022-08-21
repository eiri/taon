package taon

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/mattn/go-tty"
)

// Header is an alias for slice of strings used to define headers
type Header []string

// makeHeader creates header list
func makeHeader(m map[string]int, c Columns) (Header, error) {
	header := Header{}

	if len(m) == 0 {
		return header, errors.New("Record is empty")
	}

	if len(c) > 0 {
		var tmp []string
		for _, key := range c {
			if _, ok := m[key]; ok {
				tmp = append(tmp, key)
			}
		}
		if len(tmp) == 0 {
			return header, errors.New("Can't find specified column(s)")
		}
		header = tmp
		return header, nil
	}

	for key := range m {
		header = append(header, key)
	}
	// otherwise we can't guarantee stable columns order
	sort.Strings(header)
	return header, nil
}

// makeCell converts from typed input to string representation
func makeCell(v interface{}) string {
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

func max(ints ...int) int {
	m := ints[0]
	for i := range ints {
		if ints[i] > m {
			m = ints[i]
		}
	}
	return m
}

// maxColumns returns tty's width
func maxColumns() (int, error) {
	// prefer env $COLUMNS, fail back on tty if not set
	if columns := os.Getenv("COLUMNS"); columns != "" {
		return strconv.Atoi(columns)
	}

	tty, err := tty.Open()
	if err != nil {
		return 0, err
	}
	defer tty.Close()
	_, width, err := tty.Size()
	return width, err
}
