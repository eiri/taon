package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"sort"
)

// Header is an alias for slice of strings used to define headers
type Header []string

// Rows is an alias for slice of strings' slices representing table rows
type Rows [][]string

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

	if *md {
		renderMarkdown(w, header, rows)
	} else {
		renderTable(w, header, rows)
	}
	os.Exit(exitOK)
}

func renderTable(w io.Writer, header Header, rows Rows) {
	table := tablewriter.NewWriter(w)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader(header)
	table.AppendBulk(rows)
	table.Render()
}

func renderMarkdown(w io.Writer, header Header, rows Rows) {
	table := tablewriter.NewWriter(w)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Right: true})
	table.SetCenterSeparator("|")
	table.SetHeader(header)
	table.AppendBulk(rows)
	table.Render()
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
			rec = r
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
		var row []string
		for _, key := range header {
			cell := fmt.Sprintf("%v", rec[key])
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
