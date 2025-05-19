package taon

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"sort"

	"github.com/goccy/go-json"
	"github.com/olekukonko/tablewriter"
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
	decoder := json.NewDecoder(r)
	decoder.UseNumber()

	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("error reading token: %w", err)
	}

	switch delim := token.(type) {
	case json.Delim:
		switch delim {
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

			for decoder.More() {
				// Read the opening '{'
				if _, err := decoder.Token(); err != nil {
					return fmt.Errorf("error reading opening token: %w", err)
				}
				item, err := parseObject(decoder)
				if err != nil {
					return err
				}

				list = append(list, item)
			}

			// Consume the closing ']'
			if _, err := decoder.Token(); err != nil && err != io.EOF {
				return fmt.Errorf("error reading closing token: %w", err)
			}

			t.buildArray(list)
		default:
			return fmt.Errorf("unexpected delimiter: %q", delim)
		}
	default:
		return fmt.Errorf("unexpected token (expected '{' or '['): %q", token)
	}

	table := tablewriter.NewWriter(w)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	if t.isMarkdown {
		table.SetBorders(tablewriter.Border{
			Left:   true,
			Top:    false,
			Right:  true,
			Bottom: false,
		})
		table.SetCenterSeparator("|")
	}

	table.SetHeader(t.headers)
	for _, row := range t.rows {
		table.Append(row)
	}
	table.Render()

	return nil
}

func parseObject(d *json.Decoder) (map[string]string, error) {
	obj := make(map[string]string)

	for d.More() {
		keyToken, err := d.Token()
		if err != nil {
			return nil, fmt.Errorf("error reading key token: %w", err)
		}
		key, ok := keyToken.(string)
		if !ok {
			return nil, fmt.Errorf("expected string key, got: %q", keyToken)
		}

		var val any
		if err := d.Decode(&val); err != nil {
			return nil, fmt.Errorf("error decoding value: %w", err)

		}
		flat := make(map[string]any)
		flatten(key, val, flat)
		for k, v := range flat {
			obj[k] = makeCell(v)
		}
	}

	// Read the closing '}'
	if _, err := d.Token(); err != nil {
		return nil, fmt.Errorf("error reading closing token: %w", err)
	}

	return obj, nil
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
