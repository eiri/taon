package main

import (
	"encoding/json"
	"errors"
	"fmt"
	tt "github.com/apcera/termtables"
	ff "github.com/jeremywohl/flatten"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"sort"
	"strconv"
)

// Header is an alias for slice of strings used to define headers
type Header []string

// Rows is an alias for slice of strings' slices representing table rows
type Rows [][]interface{}

const (
	exitOK = iota
	exitOpenFile
	exitParseError
	exitRenderTable
)

var (
	version = "dev"
	columns = &ColumnsValue{}
	file    *string
	md      *bool
)

func main() {
	var r io.Reader
	var w io.Writer
	var err error
	w = os.Stdout

	taon := kingpin.New("taon", "Transform JSON into ASCII table.")
	taon.Version(version)
	taon.HelpFlag.Short('h')
	s := taon.Flag("columns", "List of columns to display").
		PlaceHolder("COL1,COL2").Short('c')
	s.SetValue((*ColumnsValue)(columns))
	md = taon.Flag("markdown", "Print as markdown table").Short('m').Bool()
	file = taon.Arg("file", "File to read").ExistingFile()
	taon.Parse(os.Args[1:])

	if *file == "" {
		r = os.Stdin
	} else {
		r, err = os.Open(*file)
		if err != nil {
			taon.Errorf("Failed to open file: %s\n", err)
			os.Exit(exitOpenFile)
		}
	}

	header, rows, err := parseJSON(r)
	if err != nil {
		taon.Errorf("Failed to parse JSON: %s\n", err)
		os.Exit(exitParseError)
	}

	table := makeTable(header, rows)

	_, err = fmt.Fprint(w, table)
	if err != nil {
		taon.Errorf("Failed to render table: %s\n", err)
		os.Exit(exitRenderTable)
	}

	os.Exit(exitOK)
}

func makeTable(header Header, rows Rows) string {
	table := tt.CreateTable()
	if *md {
		table.SetModeMarkdown()
	}
	for _, h := range header {
		table.AddHeaders(h)
	}
	for _, row := range rows {
		table.AddRow(row...)
	}
	return table.Render()
}

func parseJSON(r io.Reader) (header Header, rows Rows, err error) {
	var raw interface{}
	d := json.NewDecoder(r)
	d.UseNumber()
	err = d.Decode(&raw)
	if err != nil {
		return
	}
	var records []interface{}
	switch rec := raw.(type) {
	case []interface{}:
		records = rec
	case map[string]interface{}:
		records = append(records, rec)
	default:
		err = errors.New("Unsupported JSON data structure")
		return
	}

	tableble := true
	for _, val := range records {
		var rec map[string]interface{}
		if r, ok := val.(map[string]interface{}); ok && tableble {
			rec, err = ff.Flatten(r, "", ff.DotStyle)
			if err != nil {
				break
			}
		} else {
			rec = map[string]interface{}{"value": val}
			tableble = false
		}
		if len(header) == 0 {
			header, err = makeHeader(rec)
			if err != nil {
				break
			}
		}
		var row []interface{}
		for _, key := range header {
			cell := makeCell(rec[key])
			row = append(row, cell)
		}
		rows = append(rows, row)
	}

	return
}

func makeHeader(m map[string]interface{}) (header Header, err error) {
	for key := range m {
		header = append(header, key)
	}
	sort.Strings(header)
	if len(*columns) > 0 {
		var tmp []string
		for _, key := range *columns {
			i := sort.SearchStrings(header, key)
			if i < len(header) && header[i] == key {
				tmp = append(tmp, key)
			}
		}
		header = tmp
	}
	if len(header) == 0 {
		err = errors.New("Can't find specified column(s)")
	}
	return
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
