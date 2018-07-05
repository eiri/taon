package main

import (
	"os"
	"reflect"
	"strconv"
	"testing"
)

type headertestpair struct {
	columns *ColumnsValue
	expect  Header
}

var headertests = []headertestpair{
	{&ColumnsValue{}, Header{"a", "b", "o", "r", "z"}},
	{&ColumnsValue{"r", "a"}, Header{"r", "a"}},
	{&ColumnsValue{"b", "x", "z", "p"}, Header{"b", "z"}},
}

// TestMakeHeader to ensure we are getting sorted list of strings
func TestMakeHeader(t *testing.T) {
	in := map[string]interface{}{"b": 1, "z": 2, "a": 3, "r": 4, "o": 5}
	for _, pair := range headertests {
		columns = pair.columns
		out, err := makeHeader(in)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(pair.expect, out) {
			t.Errorf("Expecting %#v, got %#v", pair.expect, out)
		}
	}
	// return error
	columns = &ColumnsValue{"none"}
	out, err := makeHeader(in)
	if err == nil {
		t.Error("Expecting error, got nil")
	}
	var expect Header
	if !reflect.DeepEqual(expect, out) {
		t.Errorf("Expecting %#v, got %#v", expect, out)
	}
}

// TestParseJSONObject to ensure we can parse JSON object
func TestParseJSONObject(t *testing.T) {
	columns = &ColumnsValue{}
	r, err := os.Open("testdata/object.json")
	if err != nil {
		t.Fatal(err)
	}
	header, rows, err := parseJSON(r)
	if err != nil {
		t.Fatal(err)
	}
	expectHeader := Header{"bool", "int", "string"}
	if !reflect.DeepEqual(expectHeader, header) {
		t.Errorf("Expecting %#v, got %#v", expectHeader, header)
	}
	expectRows := Rows{[]interface{}{"true", "42", "answer"}}
	if !reflect.DeepEqual(expectRows, rows) {
		t.Errorf("Expecting %#v, got %#v", expectRows, rows)
	}
}

// TestParseJSONArray to ensure we can parse array of JSON objects
func TestParseJSONArray(t *testing.T) {
	columns = &ColumnsValue{}
	r, err := os.Open("testdata/array.json")
	if err != nil {
		t.Fatal(err)
	}
	header, rows, err := parseJSON(r)
	if err != nil {
		t.Fatal(err)
	}
	expectHeader := Header{"#", "char"}
	if !reflect.DeepEqual(expectHeader, header) {
		t.Errorf("Expecting %#v, got %#v", expectHeader, header)
	}
	var expectRows Rows
	for i, l := range "abcde" {
		row := []interface{}{strconv.Itoa(i + 1), string(l)}
		expectRows = append(expectRows, row)
	}
	if !reflect.DeepEqual(expectRows, rows) {
		t.Errorf("Expecting %#v, got %#v", expectRows, rows)
	}
}

// TestParseMiscJSONArray to ensure we can parse arbitrary JSON array
func TestParseMiscJSONArray(t *testing.T) {
	columns = &ColumnsValue{}
	r, err := os.Open("testdata/misc-array.json")
	if err != nil {
		t.Fatal(err)
	}
	header, rows, err := parseJSON(r)
	if err != nil {
		t.Fatal(err)
	}
	expectHeader := Header{"value"}
	if !reflect.DeepEqual(expectHeader, header) {
		t.Errorf("Expecting %#v, got %#v", expectHeader, header)
	}
	expectRows := Rows{
		[]interface{}{"309"},
		[]interface{}{"true"},
		[]interface{}{"zG8dnbd1iXDHAewJ"},
		[]interface{}{"false"},
		[]interface{}{"773"},
		[]interface{}{"Og3TQltUz2eIW6ZF"},
		[]interface{}{"map[note:c#]"},
	}
	if !reflect.DeepEqual(expectRows, rows) {
		t.Errorf("Expecting %#v, got %#v", expectRows, rows)
	}
}
