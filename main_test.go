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
	in := map[string]int{"b": 1, "z": 2, "a": 3, "r": 4, "o": 5}
	for _, pair := range headertests {
		columns = pair.columns
		out := makeHeader(in)
		if !reflect.DeepEqual(pair.expect, out) {
			t.Errorf("Expecting %#v, got %#v", pair.expect, out)
		}
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
	expectRows := Rows{[]string{"true", "42", "answer"}}
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
		row := []string{strconv.Itoa(i + 1), string(l)}
		expectRows = append(expectRows, row)
	}
	if !reflect.DeepEqual(expectRows, rows) {
		t.Errorf("Expecting %#v, got %#v", expectRows, rows)
	}
}

// TestColumnsValue to complaiance to flag's Value type
func TestColumnsValue(t *testing.T) {
	var c ColumnsValue
	if c.String() != "[]" {
		t.Errorf("Expecting `[]` for zero ColumnsValue")
	}
	c.Set("zebra,alpha,comma")
	expect := ColumnsValue{"zebra", "alpha", "comma"}
	if !reflect.DeepEqual(expect, c) {
		t.Errorf("Expecting %#v, got %#v", expect, c)
	}
	if c.String() != "[zebra alpha comma]" {
		t.Errorf("Expecting `[zebra alpha comma]` for ColumnsValue")
	}

}
