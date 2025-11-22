package taon

import (
	"bytes"
	"encoding/json"
	"encoding/json/jsontext"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

var names = []string{"array", "data", "data_deep", "data_object", "long_field", "object"}

var table_parse_object_tests = []struct {
	json   string
	expect map[string]string
}{
	{
		`{"a": "b", "c": 1, "d": false}`,
		map[string]string{"a": "b", "c": "1", "d": "false"},
	},
	{
		`{"a": "b", "c": {"d": {"e": "f", "g": "h"}}}`,
		map[string]string{"a": "b", "c.d.e": "f", "c.d.g": "h"},
	},
	{
		`{"a": "b", "c": {"d": {"e": "f"}, "g": {"h": "k"}}}`,
		map[string]string{"a": "b", "c.d.e": "f", "c.g.h": "k"},
	},
	{
		`{"a": "b", "c": [1, 2]}`,
		map[string]string{"a": "b", "c.0": "1", "c.1": "2"},
	},
	{
		`{"a": "b", "c": [{"d": 1, "e": 2}]}`,
		map[string]string{"a": "b", "c.0.d": "1", "c.0.e": "2"},
	},
	{
		`{"a": {"b": 1, "c": 2}, "d": [{"e": 1}, {"f": 2}]}`,
		map[string]string{"a.b": "1", "a.c": "2", "d.0.e": "1", "d.1.f": "2"},
	},
}

// TestParseObject to confirm that we parse and flatten json objects
func TestParseObject(t *testing.T) {
	for _, tt := range table_parse_object_tests {
		dec := jsontext.NewDecoder(bytes.NewBufferString(tt.json))
		parsed, err := parseObject(dec)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(parsed, tt.expect) {
			t.Errorf("expected %s to be parsed into %#v", tt.json, tt.expect)
		}
	}
}

// FuzzParseObject checks that parseObject handles arbitrary valid JSON.
func FuzzParseObject(f *testing.F) {
	// Seed with some structured values
	f.Add(`{"a":"b","c":{"d":{"e":"f"}}}`)
	f.Add(`{"arr":[1,2,{"x":3}]}`)
	f.Add(`{}`)
	f.Add(`{"nested":[{"a":1},{"b":[2,3]}]}`)

	f.Fuzz(func(t *testing.T, input string) {
		var js any
		if err := json.Unmarshal([]byte(input), &js); err != nil {
			// skip non-JSON
			return
		}

		dec := jsontext.NewDecoder(bytes.NewBufferString(input))
		m, err := parseObject(dec)
		if err != nil {
			// The parser requires top-level object,
			// so reject non-object JSON.
			var obj map[string]any
			if json.Unmarshal([]byte(input), &obj) == nil {
				t.Errorf("parseObject failed on valid object JSON: %v", err)
			}
			return
		}

		// Validate resulting map
		for k, v := range m {
			if k == "" {
				t.Fatalf("empty key in result map")
			}
			if v == "" {
				// empty string allowed only if original JSON had "" or null → "null"
				_ = v
			}

			// Keys should be dot-separated segments, no empty internal segments like "a..b"
			if len(k) > 0 && k[0] == '.' {
				t.Fatalf("key begins with dot: %q", k)
			}
			if len(k) > 0 && k[len(k)-1] == '.' {
				t.Fatalf("key ends with dot: %q", k)
			}
			// check if key contains double dot
			for i := 0; i+1 < len(k); i++ {
				if k[i] == '.' && k[i+1] == '.' {
					t.Fatalf("key contains invalid '..': %q", k)
				}
			}
		}
	})
}

// TestRender to confirm that our output is formatted as table
func TestRender(t *testing.T) {
	for _, name := range names {
		r, err := os.Open("testdata/" + name + ".json")
		if err != nil {
			t.Fatal(err)
		}

		w := new(bytes.Buffer)
		table := NewTable()
		err = table.Render(r, w)
		if err != nil {
			t.Fatal(err)
		}
		r.Close()

		expect, err := os.ReadFile("testdata/" + name + ".txt")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(expect, w.Bytes()) {
			t.Errorf("for %s.txt expecting:\n%sgot:\n%s", name, expect, w.Bytes())
		}
	}
}

// BenchmarkRender bench table render
func BenchmarkRender(b *testing.B) {
	for _, rows := range []int{1, 10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d_rows", rows), func(b *testing.B) {

			prep := bytes.NewBufferString("[")
			for i := range rows {
				fmt.Fprintf(prep, `{"num": "%06d", "field": "%06X"}`, i+1, (i+1)*1000000)
				if i < rows-1 {
					prep.WriteRune(',')
				}
			}
			prep.WriteRune(']')
			data := prep.Bytes()
			table := NewTable()

			b.ResetTimer()

			for b.Loop() {
				r := bytes.NewReader(data)
				w := new(bytes.Buffer)
				if err := table.Render(r, w); err != nil {
					b.Fatal(err)
				}
			}
		})
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
		table := NewTable()
		table.SetModeMarkdown()

		err = table.Render(r, w)
		if err != nil {
			t.Fatal(err)
		}
		r.Close()

		expect, err := os.ReadFile("testdata/" + name + ".md")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(expect, w.Bytes()) {
			t.Errorf("for %s.md expecting:\n%sgot:\n%s", name, expect, w.Bytes())
		}
	}
}

var table_columns_tests = []struct {
	name    string
	columns Columns
}{
	{"data", Columns{"seq", "name", "word"}},
	{"data_deep", Columns{"key", "value.rev", "doc.name"}},
	{"data_object", Columns{"key", "value.rev", "doc.name"}},
}

// TestRenderColumns to confirm that our output table has reduced columns
func TestRenderColumns(t *testing.T) {
	for _, tt := range table_columns_tests {
		r, err := os.Open("testdata/" + tt.name + ".json")
		if err != nil {
			t.Fatal(err)
		}

		w := new(bytes.Buffer)
		table := NewTable()
		table.SetColumns(tt.columns)

		err = table.Render(r, w)
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

	table := NewTable()
	err = table.Render(r, io.Discard)
	if err == nil {
		t.Error("Expecting error, got nil")
	}
}
