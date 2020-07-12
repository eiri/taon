package main

import (
	"fmt"
	"io"
	"os"

	"github.com/eiri/taon/pkg/taon"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	exitOK = iota
	exitOpenFile
	exitParseError
	exitRenderTable
)

var (
	version = "dev"
	columns = &taon.ColumnsValue{}
	file    *string
	md      *bool
)

func main() {
	var r io.Reader
	var w io.Writer
	var err error
	w = os.Stdout

	app := kingpin.New("taon", "Transform JSON into ASCII table.")
	app.Version(version)
	app.HelpFlag.Short('h')
	s := app.Flag("columns", "List of columns to display").
		PlaceHolder("COL1,COL2").Short('c')
	s.SetValue((*taon.ColumnsValue)(columns))
	md = app.Flag("markdown", "Print as markdown table").Short('m').Bool()
	file = app.Arg("file", "File to read").ExistingFile()
	app.Parse(os.Args[1:])

	if *file == "" {
		r = os.Stdin
	} else {
		r, err = os.Open(*file)
		if err != nil {
			app.Errorf("Failed to open file: %s\n", err)
			os.Exit(exitOpenFile)
		}
	}

	t := taon.NewTable()

	if *md {
		t.SetModeMarkdown()
	}

	if len(*columns) > 0 {
		t.SetColumns(*columns)
	}

	table, err := t.Render(r)
	if err != nil {
		app.Errorf("Failed to render table: %s\n", err)
		os.Exit(exitParseError)
	}

	_, err = fmt.Fprint(w, table)
	if err != nil {
		os.Exit(exitRenderTable)
	}

	os.Exit(exitOK)
}
