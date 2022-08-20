package taon

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var names = []string{"array", "data", "data_deep", "long_field", "object"}

// TestRender to confirm that our output is formatted as table
func TestRender(t *testing.T) {
	for _, name := range names {
		r, err := os.Open("testdata/" + name + ".json")
		if err != nil {
			t.Fatal(err)
		}

		w := new(bytes.Buffer)
		table := NewTable(r, w)
		err = table.Render()
		if err != nil {
			t.Fatal(err)
		}
		r.Close()

		expect, err := os.ReadFile("testdata/" + name + ".txt")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(expect, w.Bytes()) {
			t.Errorf("Expecting:\n%sgot:\n%s", expect, w.Bytes())
		}
	}
}

// TestRenderMarkdown to confirm that our output is formatted as markdown
func TestRenderMarkdown(t *testing.T) {
	for _, name := range names {
		r, err := os.Open("testdata/" + name + ".json")
		if err != nil {
			t.Fatal(err)
		}

		w := new(bytes.Buffer)
		table := NewTable(r, w)
		table.SetModeMarkdown()

		err = table.Render()
		if err != nil {
			t.Fatal(err)
		}
		r.Close()

		expect, err := os.ReadFile("testdata/" + name + ".md")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(expect, w.Bytes()) {
			t.Errorf("Expecting:\n%sgot:\n%s", expect, w.Bytes())
		}
	}
}

var table_columns_tests = []struct {
	name    string
	columns Columns
}{
	{"data", Columns{"seq", "name", "word"}},
	{"data_deep", Columns{"key", "value.rev", "doc.name"}},
}

// TestRenderColumns to confirm that our output table has reduced columns
func TestRenderColumns(t *testing.T) {
	for _, tt := range table_columns_tests {
		r, err := os.Open("testdata/" + tt.name + ".json")
		if err != nil {
			t.Fatal(err)
		}

		w := new(bytes.Buffer)
		table := NewTable(r, w)
		table.SetColumns(tt.columns)

		err = table.Render()
		if err != nil {
			t.Fatal(err)
		}
		r.Close()

		expect, err := os.ReadFile("testdata/" + tt.name + "_columns.txt")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(expect, w.Bytes()) {
			t.Errorf("Expecting:\n%sgot:\n%s", expect, w.Bytes())
		}
	}
}

// TestErrorMiscJSONArray to ensure we are returning error on arbitrary JSON array
func TestErrorMiscJSONArray(t *testing.T) {
	r, err := os.Open("testdata/misc_array.json")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	table := NewTable(r, ioutil.Discard)
	err = table.Render()
	if err == nil {
		t.Error("Expecting error, got nil")
	}
}
