package taon

import (
	"encoding/json"
	"errors"
	"io"

	ff "github.com/jeremywohl/flatten/v2"
	"github.com/olekukonko/tablewriter"
)

// Table datastructure represents ascii table
type Table struct {
	columns Columns
	reader  io.Reader
	tw      *tablewriter.Table
}

// NewTable initializes and returns new Table
func NewTable(r io.Reader, w io.Writer) *Table {
	table := tablewriter.NewWriter(w)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	return &Table{
		columns: Columns{},
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

	lookup := make(map[string]bool)
	if len(t.columns) > 0 {
		for _, k := range t.columns {
			lookup[k] = true
		}
	}
	columns := make(map[string]int)
	rows := make([]map[string]string, 0)
	for _, val := range records {
		r, ok := val.(map[string]interface{})
		if !ok {
			return errors.New("Unsupported JSON data structure")
		}

		rec, err := ff.Flatten(r, "", ff.DotStyle)
		if err != nil {
			return err
		}

		row := make(map[string]string)
		for key, value := range rec {
			_, ok := lookup[key]
			if !ok && len(lookup) > 0 {
				continue
			}
			if _, ok := columns[key]; !ok {
				columns[key] = 0
			}
			row[key] = makeCell(value)
			columns[key] = max(columns[key], len(key), len(row[key]))
		}
		rows = append(rows, row)
	}

	header, err := makeHeader(columns, t.columns)
	if err != nil {
		return err
	}
	t.tw.SetHeader(header)

	// calc length of each column
	maxColumns, err := maxColumns()
	if err != nil {
		return err
	}
	maxlen := (maxColumns - 3*len(header) - 1) / len(header)
	// margin is a global length cols can grow to in place of short columns
	margin := 0
	for _, l := range columns {
		if l < maxlen {
			margin += maxlen - l
		}
	}

	for _, k := range header {
		l := columns[k]
		if l > maxlen {
			need := l - maxlen
			// give all margin to this column
			if need > margin {
				columns[k] = maxlen + margin
				margin = 0
			} else {
				columns[k] = maxlen + need
				margin -= need
			}
		}
	}

	// build rows, cut cells to each col width
	for _, row := range rows {
		r := make([]string, 0)
		for _, key := range header {
			l := columns[key]
			if l > 3 && len(row[key]) > l {
				row[key] = row[key][:l-3] + "..."
			}
			r = append(r, row[key])
		}
		t.tw.Append(r)
	}

	t.tw.Render()

	return nil
}
