package taon

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
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

func (t *Table) JsontextParser(r io.Reader) error {
	dec := jsontext.NewDecoder(r)
	for {
		// Read a token from the input.
		tok, err := dec.ReadToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fmt.Printf("kind: %c tok: %s\n", tok.Kind(), tok.String())
	}
	return nil
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

func parseObject(d *jsontext.Decoder) (map[string]string, error) {
	obj := make(map[string]string)

	// discard opening token
	if _, err := d.ReadToken(); err != nil {
		return nil, err
	}

	for d.PeekKind() != '}' {
		keyToken, err := d.ReadToken()
		if err != nil {
			return nil, fmt.Errorf("error reading key token: %w", err)
		}
		key := keyToken.String()

		var val any
		err = json.UnmarshalDecode(d, &val,
			json.WithUnmarshalers(
				json.UnmarshalFromFunc(func(dec *jsontext.Decoder, val *any) error {
					if dec.PeekKind() == '0' {
						*val = jsontext.Value(nil)
					}
					return json.SkipFunc
				}),
			))
		if err != nil {
			return nil, fmt.Errorf("error decoding value: %w", err)
		}
		flat := make(map[string]any)
		flatten(key, val, flat)
		for k, v := range flat {
			obj[k] = makeCell(v)
		}
	}

	// discard closing token
	if _, err := d.ReadToken(); err != nil {
		return nil, err
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
