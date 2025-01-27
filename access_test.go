package fluffyjson_test

import (
	"fmt"
	"testing"

	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func TestAccess(t *testing.T) {
	testcases := map[string]struct {
		target   string
		accessor fluffyjson.Accessor
		expect   fluffyjson.JsonValue
		err      error
	}{
		"key access": {
			target:   `{"hello":"world"}`,
			accessor: fluffyjson.KeyAccess("hello"),
			expect:   HelperCastString(t, "world"),
			err:      nil,
		},
		"index access": {
			target:   `["hello", "world"]`,
			accessor: fluffyjson.IndexAccess(1),
			expect:   HelperCastString(t, "world"),
			err:      nil,
		},
		"invalid key access": {
			target:   `["hello", "world"]`,
			accessor: fluffyjson.KeyAccess("hello"),
			expect:   nil,
			err: fluffyjson.ErrAccess{
				Accessor: fmt.Sprintf("%T", fluffyjson.KeyAccess("hello")),
				Expected: fluffyjson.OBJECT,
				Actual:   fluffyjson.ARRAY,
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			value := HelperUnmarshalValue(t, tc.target)
			actual, err := value.Access(tc.accessor)
			HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
		})
	}
}

func TestSliceAccess(t *testing.T) {
	testcases := map[string]struct {
		target   string
		accessor fluffyjson.SliceAccessor
		expect   []fluffyjson.JsonValue
		err      error
	}{
		"slice access": {
			target:   `["one", "two", "three"]`,
			accessor: fluffyjson.SliceAccess{Start: 1, End: 3},
			expect: []fluffyjson.JsonValue{
				HelperCastString(t, "two"),
				HelperCastString(t, "three"),
			},
			err: nil,
		},
		"invalid slice access": {
			target:   `{"hello":"world"}`,
			accessor: fluffyjson.SliceAccess{Start: 0, End: 2},
			expect:   nil,
			err: fluffyjson.ErrAccess{
				Accessor: fmt.Sprintf("%T", fluffyjson.SliceAccess{}),
				Expected: fluffyjson.ARRAY,
				Actual:   fluffyjson.OBJECT,
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			value := HelperUnmarshalValue(t, tc.target)
			actual, err := value.Slice(tc.accessor)
			HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
		})
	}
}

func TestPointer(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		testcases := map[string]struct {
			target  string
			pointer fluffyjson.Pointer
			expect  fluffyjson.JsonValue
			err     error
		}{
			"root": {
				target:  `{"hello":"world"}`,
				pointer: fluffyjson.ParsePointer("/"),
				expect:  &fluffyjson.Object{"hello": HelperCastString(t, "world")},
				err:     nil,
			},
			"slice access": {
				target:  `{"number": ["zero", "one", "two"]}`,
				pointer: fluffyjson.ParsePointer("/number/1"),
				expect:  HelperCastString(t, "one"),
				err:     nil,
			},
			"integer like map key": {
				target:  `{"0": "zero", "1": "one", "2": "two"}`,
				pointer: fluffyjson.ParsePointer("/0"),
				expect:  HelperCastString(t, "zero"),
				err:     nil,
			},
			"escape": {
				target:  `{"a/b~c~1": "success"}`,
				pointer: fluffyjson.ParsePointer("/a~1b~0c~01"),
				expect:  HelperCastString(t, "success"),
				err:     nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.Access(tc.pointer)
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("to string", func(t *testing.T) {
		testcases := map[string]struct {
			pointer fluffyjson.Pointer
			expect  string
		}{
			"root": {
				pointer: nil,
				expect:  "/",
			},
			"slice access": {
				pointer: fluffyjson.Pointer{fluffyjson.KeyAccess("number"), fluffyjson.IndexAccess(1)},
				expect:  "/number/1",
			},
			"escape": {
				pointer: fluffyjson.Pointer{fluffyjson.KeyAccess("a/b~c~1")},
				expect:  "/a~1b~0c~01",
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				actual := tc.pointer.String()
				HelperFatalEvaluate(t, tc.expect, actual)
			})
		}
	})
}

func TestAccessVariadic(t *testing.T) {
	t.Run("variadic parameter", func(t *testing.T) {
		target := `{"number": ["zero", "one", "two"]}`
		value := HelperUnmarshalValue(t, target)

		actual, err := value.Access(
			fluffyjson.KeyAccess("number"),
			fluffyjson.IndexAccess(2),
		)
		if err != nil {
			t.Fatal(err)
		}

		var expect fluffyjson.JsonValue = HelperCastString(t, "two")
		HelperFatalEvaluate(t, expect, actual)
	})
}
