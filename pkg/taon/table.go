package taon

import (
	"errors"
	"io"
	"sort"

	"github.com/goccy/go-json"
	ff "github.com/jeremywohl/flatten/v2"
	"github.com/olekukonko/tablewriter"
)

const (
	keysTitle   = "keys"
	valuesTitle = "values"
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
	switch records := raw.(type) {
	case []interface{}:
		return t.renderArray(records)
	case map[string]interface{}:
		return t.renderObject(records)
	}

	return errors.New("Unsupported JSON data structure")
}

// renderObject generates ascii table for object data
func (t *Table) renderObject(records map[string]interface{}) error {
	flatten, err := ff.Flatten(records, "", ff.DotStyle)
	if err != nil {
		return err
	}
	keys := t.columns
	if len(keys) == 0 {
		for key := range flatten {
			keys = append(keys, key)
		}
		// otherwise we can't guarantee stable rows order
		sort.Strings(keys)
	}

	lookup := make(map[string]bool)
	if len(t.columns) > 0 {
		for _, k := range t.columns {
			lookup[k] = true
		}
	}
	columns := map[string]int{keysTitle: 0, valuesTitle: 0}
	rows := make([][]string, 0)
	for _, key := range keys {
		_, ok := lookup[key]
		if !ok && len(lookup) > 0 {
			continue
		}
		value := flatten[key]

		keyCell := makeCell(key)
		columns[keysTitle] = max(columns[keysTitle], len(keyCell))

		valueCell := makeCell(value)
		columns[valuesTitle] = max(columns[valuesTitle], len(valueCell))

		rows = append(rows, []string{keyCell, valueCell})
	}

	header := []string{keysTitle, valuesTitle}
	t.tw.SetHeader(header)

	columns, err = calcColumnsWidth(columns, header)
	if err != nil {
		return err
	}

	// build rows, cut cells to each col width
	for _, row := range rows {
		k, v := row[0], row[1]
		kl, vl := columns[keysTitle], columns[valuesTitle]
		if kl > 3 && len(k) > kl {
			row[0] = k[:kl-3] + "..."
		}
		if vl > 3 && len(v) > vl {
			row[1] = v[:vl-3] + "..."
		}
		t.tw.Append(row)
	}

	t.tw.Render()

	return nil
}

// renderArray generates ascii table for array data
func (t *Table) renderArray(records []interface{}) error {
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

	columns, err = calcColumnsWidth(columns, header)
	if err != nil {
		return err
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

func calcColumnsWidth(columns map[string]int, header Header) (map[string]int, error) {
	// calc length of each column
	maxColumns, err := maxColumns()
	if err != nil {
		return columns, err
	}
	// margin is a global length cols can grow to in place of short columns
	maxlen, margin := (maxColumns-3*len(header)-1)/len(header), 0
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

	return columns, nil
}
