package taon

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"testing"
)

var headers_tests = []struct {
	columns Columns
	expect  Header
}{
	{Columns{}, Header{"a", "b", "o", "r", "z"}},
	{Columns{"r", "a"}, Header{"r", "a"}},
	{Columns{"b", "x", "z", "p"}, Header{"b", "z"}},
}

// TestMakeHeader to ensure we are getting sorted list of strings
func TestMakeHeader(t *testing.T) {
	in := map[string]interface{}{"b": 1, "z": 2, "a": 3, "r": 4, "o": 5}
	for _, tt := range headers_tests {
		table := NewTable(os.Stdin, ioutil.Discard)
		table.SetColumns(tt.columns)
		err := table.makeHeader(in)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tt.expect, table.header) {
			t.Errorf("Expecting %#v, got %#v", tt.expect, table.header)
		}
	}

	// return error
	table := NewTable(os.Stdin, ioutil.Discard)
	table.SetColumns(Columns{"none"})
	err := table.makeHeader(in)
	if err == nil {
		t.Error("Expecting error, got nil")
	}
	expect := Header{}
	if !reflect.DeepEqual(expect, table.header) {
		t.Errorf("Expecting %#v, got %#v", expect, table.header)
	}
}

// TestParseJSONObject to ensure we can parse JSON object
func TestParseJSONObject(t *testing.T) {
	r, err := os.Open("testdata/object.json")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	table := NewTable(r, ioutil.Discard)
	err = table.Render()
	if err != nil {
		t.Fatal(err)
	}
	expectHeader := Header{"bool", "int", "string"}
	if !reflect.DeepEqual(expectHeader, table.header) {
		t.Errorf("Expecting %#v, got %#v", expectHeader, table.header)
	}
	expectRows := Rows{[]string{"true", "42", "answer"}}
	if !reflect.DeepEqual(expectRows, table.rows) {
		t.Errorf("Expecting %#v, got %#v", expectRows, table.rows)
	}
}

// TestParseJSONArray to ensure we can parse array of JSON objects
func TestParseJSONArray(t *testing.T) {
	r, err := os.Open("testdata/array.json")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	table := NewTable(r, ioutil.Discard)
	err = table.Render()
	if err != nil {
		t.Fatal(err)
	}
	expectHeader := Header{"#", "char"}
	if !reflect.DeepEqual(expectHeader, table.header) {
		t.Errorf("Expecting %#v, got %#v", expectHeader, table.header)
	}
	var expectRows Rows
	for i, l := range "abcde" {
		row := []string{strconv.Itoa(i + 1), string(l)}
		expectRows = append(expectRows, row)
	}
	if !reflect.DeepEqual(expectRows, table.rows) {
		t.Errorf("Expecting %#v, got %#v", expectRows, table.rows)
	}
}

// TestParseMiscJSONArray to ensure we can parse arbitrary JSON array
func TestParseMiscJSONArray(t *testing.T) {
	r, err := os.Open("testdata/misc_array.json")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	table := NewTable(r, ioutil.Discard)
	err = table.Render()
	if err != nil {
		t.Fatal(err)
	}
	expectHeader := Header{"value"}
	if !reflect.DeepEqual(expectHeader, table.header) {
		t.Errorf("Expecting %#v, got %#v", expectHeader, table.header)
	}
	expectRows := Rows{
		[]string{"309"},
		[]string{"true"},
		[]string{"zG8dnbd1iXDHAewJ"},
		[]string{"false"},
		[]string{"773"},
		[]string{"Og3TQltUz2eIW6ZF"},
		[]string{"map[note:c#]"},
	}
	if !reflect.DeepEqual(expectRows, table.rows) {
		t.Errorf("Expecting %#v, got %#v", expectRows, table.rows)
	}
}

var names = []string{"array", "data", "data_deep", "long_field", "misc_array", "object"}

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
