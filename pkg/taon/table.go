package taon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"

	ff "github.com/jeremywohl/flatten"
	tt "github.com/scylladb/termtables"
)

// Table datastructure represents ascii table
type Table struct {
	md      bool
	columns ColumnsValue
	header  Header
	rows    Rows
}

// Header is an alias for slice of strings used to define headers
type Header []string

// Rows is an alias for slice of strings' slices representing table rows
type Rows [][]interface{}

// NewTable initializes and returns new Table
func NewTable() *Table {
	return &Table{
		md:      false,
		columns: ColumnsValue{},
		header:  Header{},
		rows:    Rows{},
	}
}

// SetModeMarkdown to draw table as MD table
func (t *Table) SetModeMarkdown() {
	t.md = true
}

// SetColumns to restrict output to only given columns
func (t *Table) SetColumns(cv ColumnsValue) {
	t.columns = cv
}

// Render generates ascii table
func (t *Table) Render(r io.Reader) (string, error) {
	var raw interface{}
	d := json.NewDecoder(r)
	d.UseNumber()
	err := d.Decode(&raw)
	if err != nil {
		return "", err
	}
	var records []interface{}
	switch rec := raw.(type) {
	case []interface{}:
		records = rec
	case map[string]interface{}:
		records = append(records, rec)
	default:
		return "", errors.New("Unsupported JSON data structure")
	}

	tableble := true
	var ruler []int
	for _, val := range records {
		var rec map[string]interface{}
		if r, ok := val.(map[string]interface{}); ok && tableble {
			rec, err = ff.Flatten(r, "", ff.DotStyle)
			if err != nil {
				return "", err
			}
		} else {
			rec = map[string]interface{}{"value": val}
			tableble = false
		}
		if len(t.header) == 0 {
			err = t.makeHeader(rec)
			if err != nil {
				return "", err
			}
			ruler = make([]int, len(t.header))
		}
		var row []interface{}
		for i, key := range t.header {
			cell := makeCell(rec[key])
			row = append(row, cell)
			if ruler[i] < len(key) {
				ruler[i] = len(key)
			}
			if ruler[i] < len(cell) {
				ruler[i] = len(cell)
			}
		}
		t.rows = append(t.rows, row)
	}

	maxlen, margin := (tt.MaxColumns-3*len(t.header)-1)/len(t.header), 0
	for _, r := range ruler {
		if r < maxlen {
			margin += maxlen - r
		}
	}

	for i, r := range ruler {
		if r > maxlen {
			need := r - maxlen
			if need > margin {
				ruler[i] = maxlen + margin
				margin = 0
			} else {
				ruler[i] = maxlen + need
				margin -= need
			}
		}
	}

	for _, row := range t.rows {
		for i, cell := range row {
			ml := ruler[i]
			c := cell.(string)
			if len(c) > ml {
				row[i] = c[:ml-3] + "..."
			}
		}
	}

	table := tt.CreateTable()
	if t.md {
		table.SetModeMarkdown()
	}
	for _, h := range t.header {
		table.AddHeaders(h)
	}
	for _, row := range t.rows {
		table.AddRow(row...)
	}
	return table.Render(), nil
}

func (t *Table) makeHeader(m map[string]interface{}) error {
	for key := range m {
		t.header = append(t.header, key)
	}
	sort.Strings(t.header)
	if len(t.columns) > 0 {
		var tmp []string
		for _, key := range t.columns {
			i := sort.SearchStrings(t.header, key)
			if i < len(t.header) && t.header[i] == key {
				tmp = append(tmp, key)
			}
		}
		t.header = tmp
	}
	if len(t.header) == 0 {
		return errors.New("Can't find specified column(s)")
	}
	return nil
}

// makeCell is taken from termtable's cell.renderValue
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
