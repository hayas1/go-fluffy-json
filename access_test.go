package fluffyjson_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func TestAccess(t *testing.T) {
	testcases := []struct {
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
			expect:   fluffyjson.ForceString("world"),
			err:      nil,
		},
		{
			name:     "index access",
			target:   `["hello", "world"]`,
			accessor: fluffyjson.IndexAccess(1),
			expect:   fluffyjson.ForceString("world"),
			err:      nil,
		},
		{
			name:     "invalid key access",
			target:   `["hello", "world"]`,
			accessor: fluffyjson.KeyAccess("hello"),
			expect:   nil,
			err: fluffyjson.ErrAccess{
				Accessor: fmt.Sprintf("%T", fluffyjson.KeyAccess("hello")),
				Expect:   fluffyjson.OBJECT,
				Actual:   fluffyjson.ARRAY,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.RootValue
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			actual, err := value.Access(tc.accessor)
			if !errors.Is(err, tc.err) {
				t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
			} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSliceAccess(t *testing.T) {
	testcases := []struct {
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
				fluffyjson.ForceString("two"),
				fluffyjson.ForceString("three"),
			},
			err: nil,
		},
		{
			name:     "invalid slice access",
			target:   `{"hello":"world"}`,
			accessor: fluffyjson.SliceAccess{Start: 0, End: 2},
			expect:   nil,
			err: fluffyjson.ErrAccess{
				Accessor: fmt.Sprintf("%T", fluffyjson.SliceAccess{}),
				Expect:   fluffyjson.ARRAY,
				Actual:   fluffyjson.OBJECT,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.RootValue
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			actual, err := value.Slice(tc.accessor)
			if !errors.Is(err, tc.err) {
				t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
			} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPointer(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		testcases := []struct {
			name    string
			target  string
			pointer fluffyjson.Pointer
			expect  fluffyjson.JsonValue
			err     error
		}{
			{
				name:    "root",
				target:  `{"hello":"world"}`,
				pointer: fluffyjson.ParsePointer("/"),
				expect:  &fluffyjson.Object{"hello": fluffyjson.ForceString("world")},
				err:     nil,
			},
			{
				name:    "slice access",
				target:  `{"number": ["zero", "one", "two"]}`,
				pointer: fluffyjson.ParsePointer("/number/1"),
				expect:  fluffyjson.ForceString("one"),
				err:     nil,
			},
			{
				name:    "integer like map key",
				target:  `{"0": "zero", "1": "one", "2": "two"}`,
				pointer: fluffyjson.ParsePointer("/0"),
				expect:  fluffyjson.ForceString("zero"),
				err:     nil,
			},
			{
				name:    "escape",
				target:  `{"a/b~c~1": "success"}`,
				pointer: fluffyjson.ParsePointer("/a~1b~0c~01"),
				expect:  fluffyjson.ForceString("success"),
				err:     nil,
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				var value fluffyjson.RootValue
				if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
					t.Fatal(err)
				}

				actual, err := value.Access(tc.pointer)
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("to string", func(t *testing.T) {
		testcases := []struct {
			name    string
			pointer fluffyjson.Pointer
			expect  string
		}{
			{
				name:    "root",
				pointer: nil,
				expect:  "/",
			},
			{
				name:    "slice access",
				pointer: fluffyjson.Pointer{fluffyjson.KeyAccess("number"), fluffyjson.IndexAccess(1)},
				expect:  "/number/1",
			},
			{
				name:    "escape",
				pointer: fluffyjson.Pointer{fluffyjson.KeyAccess("a/b~c~1")},
				expect:  "/a~1b~0c~01",
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				actual := tc.pointer.String()
				if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

}

func TestAccessVariadic(t *testing.T) {
	t.Run("variadic parameter", func(t *testing.T) {
		target := `{"number": ["zero", "one", "two"]}`
		var value fluffyjson.RootValue
		if err := json.Unmarshal([]byte(target), &value); err != nil {
			t.Fatal(err)
		}

		two, err := value.Access(
			fluffyjson.KeyAccess("number"),
			fluffyjson.IndexAccess(2),
		)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(fluffyjson.ForceString("two"), two); diff != "" {
			t.Fatal(diff)
		}
	})
}
