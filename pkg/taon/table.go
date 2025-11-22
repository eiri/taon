package taon

import (
	"encoding/json/jsontext"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"sort"

	"github.com/alexeyco/simpletable"
)

const (
	titleKeys = "keys"
	titleVals = "values"
)

// Table datastructure represents ascii table
type Table struct {
	columns    Columns // FIXME! gibve me a better name
	isMarkdown bool
	headers    []string
	rows       [][]string
}

// NewTable initializes and returns new Table
func NewTable() *Table {
	return &Table{
		columns: Columns{},
	}
}

// SetModeMarkdown to draw table as MD table
func (t *Table) SetModeMarkdown() {
	t.isMarkdown = true
}

// SetColumns to restrict output to only given columns
func (t *Table) SetColumns(c Columns) {
	t.columns = c
}

// Render generates ascii table
func (t *Table) Render(r io.Reader, w io.Writer) error {
	decoder := jsontext.NewDecoder(r)

	switch k := decoder.PeekKind(); k {
	case '{':
		// It's an object
		obj, err := parseObject(decoder)
		if err != nil {
			return err
		}

		t.buildObject(obj)
	case '[':
		// It's an array
		list := make([]any, 0)

		// discard opening token
		if _, err := decoder.ReadToken(); err != nil {
			return err
		}

		for decoder.PeekKind() != ']' {
			item, err := parseObject(decoder)
			if err != nil {
				return err
			}

			list = append(list, item)
		}

		// discard closing token
		if _, err := decoder.ReadToken(); err != nil {
			return err
		}

		t.buildArray(list)
	default:
		return fmt.Errorf("unexpected token, expected '{' or '[', but got %q", k)
	}

	table := simpletable.New()

	for _, header := range t.headers {
		table.Header.Cells = append(table.Header.Cells, &simpletable.Cell{Text: header})
	}

	for _, row := range t.rows {
		r := make([]*simpletable.Cell, 0)
		for _, cell := range row {
			r = append(r, &simpletable.Cell{Text: cell})
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	if t.isMarkdown {
		table.SetStyle(simpletable.StyleMarkdown)
	}

	fmt.Fprintln(w, table.String())

	return nil
}

// parseObject parses JSON from decoder into map[string]string
func parseObject(d *jsontext.Decoder) (map[string]string, error) {
	m := make(map[string]string)
	err := parseInternal(d, m, "")
	return m, err
}

// parseInternal recursively parses JSON tokens into map[string]string.
// prefix is the dotted key path ("a.b.c").
func parseInternal(d *jsontext.Decoder, out map[string]string, prefix string) error {
	tok, err := d.ReadToken()
	if err != nil {
		return err
	}

	switch tok.Kind() {
	case '{':
		for {
			tok, err := d.ReadToken()
			if err != nil {
				return err
			}
			if tok.Kind() == '}' {
				return nil
			}

			key := tok.String()
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}

			if err := parseInternal(d, out, fullKey); err != nil {
				return err
			}
		}
	case '[':
		idx := 0
		for {
			tok := d.PeekKind()
			if tok == ']' {
				_, _ = d.ReadToken() // consume ]
				return nil
			}

			itemPrefix := fmt.Sprintf("%s.%d", prefix, idx)
			if prefix == "" {
				itemPrefix = fmt.Sprintf("%d", idx)
			}

			if err := parseInternal(d, out, itemPrefix); err != nil {
				return err
			}
			idx++
		}
	default:
		if prefix == "" {
			return errors.New("top-level value must be an object")
		}

		out[prefix] = tok.String()
		return nil
	}
}

// renderObject generates ascii table for object data
func (t *Table) buildObject(records map[string]string) error {
	keys := t.columns
	if len(keys) == 0 {
		for key := range records {
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

	columns := []int{len(titleKeys), len(titleVals)}
	rows := make([][]string, 0)
	for _, key := range keys {
		if len(lookup) > 0 {
			if _, ok := lookup[key]; !ok {
				continue
			}
		}
		value := records[key]

		columns[0] = max(columns[0], len(key))
		columns[1] = max(columns[1], len(value))
		rows = append(rows, []string{key, value})
	}

	allocated := AllocateColumnWidths(columns)

	for id, row := range rows {
		for idx, v := range row {
			vl := allocated[idx]
			if vl > 3 && len(v) > vl {
				rows[id][idx] = v[:vl-3] + "..."
			}
		}
	}

	t.headers = []string{titleKeys, titleVals}
	t.rows = rows

	return nil
}

// buildArray converts array of values into table struct
func (t *Table) buildArray(records []any) error {
	// init headers
	headers := make([]string, 0)
	if len(t.columns) > 0 {
		headers = t.columns
	} else {
		first, ok := records[0].(map[string]string)
		if !ok {
			return errors.New("unsupported JSON data structure")
		}

		for key := range maps.Keys(first) {
			headers = append(headers, key)
		}
		sort.Strings(headers)
	}

	// init columns
	columns := make([]int, len(headers))
	for idx, key := range headers {
		columns[idx] = len(key)
	}

	rows := make([][]string, 0)
	for _, val := range records {
		record, ok := val.(map[string]string)
		if !ok {
			return errors.New("unsupported JSON data structure")
		}

		row := make([]string, len(headers))
		for key, value := range record {
			idx := slices.Index(headers, key)
			if idx > -1 {
				columns[idx] = max(columns[idx], len(value))
				row[idx] = value
			}
		}
		rows = append(rows, row)
	}

	allocated := AllocateColumnWidths(columns)

	for id, row := range rows {
		for idx, v := range row {
			vl := allocated[idx]
			if vl > 3 && len(v) > vl {
				rows[id][idx] = v[:vl-3] + "..."
			}
		}
	}

	t.headers = headers
	t.rows = rows

	return nil
}
