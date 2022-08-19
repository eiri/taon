package taon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"

	ff "github.com/jeremywohl/flatten/v2"
	"github.com/mattn/go-tty"
	"github.com/olekukonko/tablewriter"
)

// Table datastructure represents ascii table
type Table struct {
	columns Columns
	header  Header
	rows    Rows
	reader  io.Reader
	tw      *tablewriter.Table
}

// Header is an alias for slice of strings used to define headers
type Header []string

// Rows is an alias for slice of strings' slices representing table rows
type Rows [][]string

// NewTable initializes and returns new Table
func NewTable(r io.Reader, w io.Writer) *Table {
	table := tablewriter.NewWriter(w)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	return &Table{
		columns: Columns{},
		header:  Header{},
		rows:    Rows{},
		reader:  r,
		tw:      table,
	}
}

// SetModeMarkdown to draw table as MD table
func (t *Table) SetModeMarkdown() {
	t.tw.SetBorders(tablewriter.Border{
		Left:   true,
		Top:    false,
		Right:  true,
		Bottom: false,
	})
	t.tw.SetCenterSeparator("|")
}

// SetColumns to restrict output to only given columns
func (t *Table) SetColumns(c Columns) {
	t.columns = c
}

// Render generates ascii table
func (t *Table) Render() error {
	var raw interface{}
	d := json.NewDecoder(t.reader)
	d.UseNumber()
	err := d.Decode(&raw)
	if err != nil {
		return err
	}
	var records []interface{}
	switch rec := raw.(type) {
	case []interface{}:
		records = rec
	case map[string]interface{}:
		records = append(records, rec)
	default:
		return errors.New("Unsupported JSON data structure")
	}

	tableble := true
	var ruler []int
	for _, val := range records {
		var rec map[string]interface{}
		if r, ok := val.(map[string]interface{}); ok && tableble {
			rec, err = ff.Flatten(r, "", ff.DotStyle)
			if err != nil {
				return err
			}
		} else {
			rec = map[string]interface{}{"value": val}
			tableble = false
		}

		if len(t.header) == 0 {
			err = t.makeHeader(rec)
			if err != nil {
				return err
			}
			ruler = make([]int, len(t.header))
		}

		var row []string
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

	maxColumns, err := maxColumns()
	if err != nil {
		return err
	}
	maxlen, margin := (maxColumns-3*len(t.header)-1)/len(t.header), 0
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
			if ml > 3 && len(cell) > ml {
				row[i] = cell[:ml-3] + "..."
			}
		}
	}

	t.tw.SetHeader(t.header)
	t.tw.AppendBulk(t.rows)
	t.tw.Render()

	return nil
}

// makeHeader creates header list
func (t *Table) makeHeader(m map[string]interface{}) error {
	if len(m) == 0 {
		return errors.New("Record is empty")
	}

	if len(t.columns) > 0 {
		var tmp []string
		for _, key := range t.columns {
			if _, ok := m[key]; ok {
				tmp = append(tmp, key)
			}
		}
		if len(tmp) == 0 {
			return errors.New("Can't find specified column(s)")
		}
		t.header = tmp
		return nil
	}

	for key := range m {
		t.header = append(t.header, key)
	}
	// otherwise we can't guarantee stable columns order
	sort.Strings(t.header)
	return nil
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

// maxColumns returns tty's width
func maxColumns() (int, error) {
	tty, err := tty.Open()
	if err != nil {
		return 0, err
	}
	defer tty.Close()
	_, width, err := tty.Size()
	return width, err
}
