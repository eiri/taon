package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/eiri/taon/pkg/taon"
)

const (
	exitOK = iota
	exitOpenFile
	exitParseError
	exitRenderTable
)

var (
	version = "dev"
	columns taon.ColumnsValue
	md      bool
	showVer bool
)

func init() {
	flag.Var(&columns, "columns", "List of columns to display")
	flag.Var(&columns, "c", "List of columns to display")
	flag.BoolVar(&md, "markdown", false, "Print as markdown table")
	flag.BoolVar(&md, "m", false, "Print as markdown table")
	flag.BoolVar(&showVer, "version", false, "Show application version")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Transform JSON into ASCII table

Usage: taon [flags] [file]

Flags:
  -c, --columns=COL1,COL2  List of columns to display
  -m, --markdown           Print as markdown table
  -h, --help               Show help
      --version            Show application version

Args:
  <file>                   Path to file to read, stdin when missing
`)
	}
}

func main() {
	flag.Parse()

	if showVer {
		fmt.Printf("taon version %s\n", version)
		os.Exit(0)
	}

	var r io.Reader
	var w io.Writer
	var err error
	w = os.Stdout

	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		r = os.Stdin
	} else {
		file := flag.Arg(0)
		r, err = os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file: %s\n", err)
			os.Exit(exitOpenFile)
		}
	}

	t := taon.NewTable()

	if md {
		t.SetModeMarkdown()
	}

	if len(columns) > 0 {
		t.SetColumns(columns)
	}

	table, err := t.Render(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to render table: %s\n", err)
		os.Exit(exitParseError)
	}

	_, err = fmt.Fprint(w, table)
	if err != nil {
		os.Exit(exitRenderTable)
	}

	os.Exit(exitOK)
}
