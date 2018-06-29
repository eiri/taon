package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"reflect"
	"sort"
)

const (
	exitOK = iota
	exitParseError
	exitOpenFile
	exitReadInput
	exitRenderTable
)

var (
	version = "dev"
	columns = &ColumnsValue{}
)

func main() {
	var r io.Reader
	var w io.Writer
	r = os.Stdin
	w = os.Stdout

	kingpin.Version(version)
	kingpin.CommandLine.HelpFlag.Short('h')
	s := kingpin.Flag("columns", "List of columns to display").
		PlaceHolder("COL1,COL2").Short('c')
	s.SetValue((*ColumnsValue)(columns))
	kingpin.Parse()

	header, rows, err := parseJSON(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(exitParseError)
	}

	renderTable(w, header, rows)
	os.Exit(exitOK)
}

func renderTable(w io.Writer, header []string, rows [][]string) {
	table := tablewriter.NewWriter(w)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader(header)
	table.AppendBulk(rows)
	table.Render()
}

func parseJSON(r io.Reader) (header []string, rows [][]string, err error) {
	var vv []interface{}
	var v interface{}
	d := json.NewDecoder(r)
	d.UseNumber()
	err = d.Decode(&v)
	if err != nil {
		return
	}
	switch v := v.(type) {
	case []interface{}:
		header = makeHeader(v[0])
		vv = v
	case map[string]interface{}:
		header = makeHeader(v)
		vv = append(vv, v)
	default:
		err = errors.New("Unsupported JSON data structure")
		return
	}
	if len(header) == 0 {
		err = errors.New("Can't find specified column(s)")
		return
	}

	for _, v := range vv {
		// we just skip none-object rows for now
		if v, ok := v.(map[string]interface{}); ok {
			var row []string
			for _, key := range header {
				val := fmt.Sprintf("%v", v[key])
				row = append(row, val)
			}
			rows = append(rows, row)
		}
	}

	return
}

func makeHeader(val interface{}) []string {
	var header []string
	r := reflect.ValueOf(val)
	for _, key := range r.MapKeys() {
		header = append(header, key.String())
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
	return header
}
