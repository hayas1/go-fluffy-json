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
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.Value
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
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.Value
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
	testcases := []struct {
		name    string
		target  string
		pointer fluffyjson.Pointer
		expect  fluffyjson.JsonValue
		err     error
	}{
		{
			name:    "slice access",
			target:  `{"number": ["zero", "one", "two"]}`,
			pointer: fluffyjson.ParsePointer("/number/1"),
			expect:  fluffyjson.ForceString("one"),
			err:     nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.Value
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			actual, err := value.Pointer(tc.pointer)
			if !errors.Is(err, tc.err) {
				t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
			} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPointerVariadic(t *testing.T) {
	t.Run("variadic parameter", func(t *testing.T) {
		target := `{"number": ["zero", "one", "two"]}`
		var value fluffyjson.Value
		if err := json.Unmarshal([]byte(target), &value); err != nil {
			t.Fatal(err)
		}

		two, err := value.Pointer(
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
