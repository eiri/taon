package taon

import (
	"os"
	"reflect"
	"strconv"
	"testing"
)

type headertestpair struct {
	columns ColumnsValue
	expect  Header
}

var headertests = []headertestpair{
	{ColumnsValue{}, Header{"a", "b", "o", "r", "z"}},
	{ColumnsValue{"r", "a"}, Header{"r", "a"}},
	{ColumnsValue{"b", "x", "z", "p"}, Header{"b", "z"}},
}

// TestMakeHeader to ensure we are getting sorted list of strings
func TestMakeHeader(t *testing.T) {
	in := map[string]interface{}{"b": 1, "z": 2, "a": 3, "r": 4, "o": 5}
	for _, pair := range headertests {
		table := NewTable()
		table.SetColumns(pair.columns)
		err := table.makeHeader(in)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(pair.expect, table.header) {
			t.Errorf("Expecting %#v, got %#v", pair.expect, table.header)
		}
	}

	// return error
	table := NewTable()
	table.SetColumns(ColumnsValue{"none"})
	err := table.makeHeader(in)
	if err == nil {
		t.Error("Expecting error, got nil")
	}
	var expect Header
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

	table := NewTable()
	_, err = table.Render(r)
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

	table := NewTable()
	_, err = table.Render(r)
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
	r, err := os.Open("testdata/misc-array.json")
	if err != nil {
		t.Fatal(err)
	}

	table := NewTable()
	_, err = table.Render(r)
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
