package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

const (
	exitOK = iota
	exitParseError
	exitOpenFile
	exitReadInput
	exitRenderTable
)

func main() {
	var r io.Reader
	var w io.Writer
	r = os.Stdin
	w = os.Stdout

	data, err := parseJSON(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitParseError)
	}

	rows := renderTable(data)
	for _, row := range rows {
		fmt.Fprintln(w, row)
	}
	os.Exit(exitOK)
}

func renderTable(d interface{}) [][]interface{} {
	var rows [][]interface{}
	v := reflect.ValueOf(d)
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		keyString := fmt.Sprintf("%v", key)
		valString := fmt.Sprintf("%v", val)
		rows = append(rows, []interface{}{keyString, valString})
	}
	return rows
}

func parseJSON(r io.Reader) (interface{}, error) {
	var v interface{}
	d := json.NewDecoder(r)
	d.UseNumber()
	err := d.Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
