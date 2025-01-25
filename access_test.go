package fluffyjson_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayash1/go-fluffy-json"
)

func TestAccess(t *testing.T) {
	testCases := []struct {
		name     string
		target   string
		accessor fluffyjson.Accessor
		expect   fluffyjson.JsonValue
		err      error
	}{
		{
			name:     "key access",
			target:   `{"hello":"world"}`,
			accessor: fluffyjson.KeyAccess("hello"),
			expect:   &[]fluffyjson.String{fluffyjson.String("world")}[0],
			err:      nil,
		},
		{
			name:     "index access",
			target:   `["hello", "world"]`,
			accessor: fluffyjson.IndexAccess(1),
			expect:   &[]fluffyjson.String{fluffyjson.String("world")}[0],
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.Value
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			actual, err := value.Access(tc.accessor)
			if err != tc.err {
				t.Fatal(err)
			} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSliceAccess(t *testing.T) {
	testCases := []struct {
		name     string
		target   string
		accessor fluffyjson.SliceAccessor
		expect   []fluffyjson.JsonValue
		err      error
	}{
		{
			name:     "slice access",
			target:   `["one", "two", "three"]`,
			accessor: fluffyjson.SliceAccess{Start: 1, End: 3},
			expect: []fluffyjson.JsonValue{
				&[]fluffyjson.String{fluffyjson.String("two")}[0],
				&[]fluffyjson.String{fluffyjson.String("three")}[0],
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.Value
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			actual, err := value.Slice(tc.accessor)
			if err != tc.err {
				t.Fatal(err)
			} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPointer(t *testing.T) {
	testCases := []struct {
		name    string
		target  string
		pointer fluffyjson.Pointer
		expect  fluffyjson.JsonValue
		err     error
	}{
		{
			name:   "slice access",
			target: `{"number": ["zero", "one", "two"]}`,
			pointer: fluffyjson.Pointer{
				fluffyjson.KeyAccess("number"),
				fluffyjson.IndexAccess(1),
			},
			expect: &[]fluffyjson.String{fluffyjson.String("one")}[0],
			err:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.Value
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			actual, err := value.Pointer(tc.pointer)
			if err != tc.err {
				t.Fatal(err)
			} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPointerVariadic(t *testing.T) {
	t.Run("variadic parameter", func(t *testing.T) {
		raw := `{"number": ["zero", "one", "two"]}`
		var value fluffyjson.Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}

		two, err := value.Pointer(
			fluffyjson.KeyAccess("number"),
			fluffyjson.IndexAccess(2),
		)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(&[]fluffyjson.String{fluffyjson.String("two")}[0], two); diff != "" {
			t.Fatal(diff)
		}
	})
}
