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
		expected fluffyjson.JsonValue
		err      error
	}{
		"key access": {
			target:   `{"hello":"world"}`,
			accessor: fluffyjson.KeyAccess("hello"),
			expected: HelperCastString(t, "world"),
			err:      nil,
		},
		"index access": {
			target:   `["hello", "world"]`,
			accessor: fluffyjson.IndexAccess(1),
			expected: HelperCastString(t, "world"),
			err:      nil,
		},
		"invalid key access": {
			target:   `["hello", "world"]`,
			accessor: fluffyjson.KeyAccess("hello"),
			expected: nil,
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
			HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
		})
	}
}

func TestSliceAccess(t *testing.T) {
	testcases := map[string]struct {
		target   string
		accessor fluffyjson.SliceAccessor
		expected []fluffyjson.JsonValue
		err      error
	}{
		"slice access": {
			target:   `["one", "two", "three"]`,
			accessor: fluffyjson.SliceAccess{Start: 1, End: 3},
			expected: []fluffyjson.JsonValue{
				HelperCastString(t, "two"),
				HelperCastString(t, "three"),
			},
			err: nil,
		},
		"invalid slice access": {
			target:   `{"hello":"world"}`,
			accessor: fluffyjson.SliceAccess{Start: 0, End: 2},
			expected: nil,
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
			HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
		})
	}
}

func TestPointer(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			expected fluffyjson.Pointer
			err      error
		}{
			"root": {
				target:   "/",
				expected: nil,
				err:      nil,
			},
		}
		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				actual, err := fluffyjson.ParsePointer(tc.target)
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})
	t.Run("access", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			pointer  string
			expected fluffyjson.JsonValue
			err      error
		}{
			"root": {
				target:   `{"hello":"world"}`,
				pointer:  "/",
				expected: &fluffyjson.Object{"hello": HelperCastString(t, "world")},
				err:      nil,
			},
			"slice access": {
				target:   `{"number": ["zero", "one", "two"]}`,
				pointer:  "/number/1",
				expected: HelperCastString(t, "one"),
				err:      nil,
			},
			"integer like map key": {
				target:   `{"0": "zero", "1": "one", "2": "two"}`,
				pointer:  "/0",
				expected: HelperCastString(t, "zero"),
				err:      nil,
			},
			"escape": {
				target:   `{"a/b~c~1": "success"}`,
				pointer:  "/a~1b~0c~01",
				expected: HelperCastString(t, "success"),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.Access(HelperFatalParsePointer(t, tc.pointer))
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("to string", func(t *testing.T) {
		testcases := map[string]struct {
			pointer  fluffyjson.Pointer
			expected string
			err      error
		}{
			"root": {
				pointer:  nil,
				expected: "/",
			},
			"slice access": {
				pointer:  fluffyjson.Pointer{fluffyjson.KeyAccess("number"), fluffyjson.IndexAccess(1)},
				expected: "/number/1",
			},
			"escape": {
				pointer:  fluffyjson.Pointer{fluffyjson.KeyAccess("a/b~c~1")},
				expected: "/a~1b~0c~01",
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				actual, err := tc.pointer.PointerString()
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
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

		var expected fluffyjson.JsonValue = HelperCastString(t, "two")
		HelperFatalEvaluate(t, expected, actual)
	})
}

func FuzzPointerRoundtrip(f *testing.F) {
	f.Add("/hoge/0/")
	f.Add("/fuga/1/~0")
	f.Fuzz(func(t *testing.T, target string) {
		pointer, err := fluffyjson.ParsePointer(target)
		if err != nil {
			return
		}
		roundtrip, err := pointer.PointerString()
		if err != nil {
			t.Fatal(err)
		}
		HelperFatalEvaluate(t, target, roundtrip)
	})
}

func HelperFatalParsePointer(t *testing.T, target string) fluffyjson.Pointer {
	t.Helper()
	pointer, err := fluffyjson.ParsePointer(target)
	if err != nil {
		t.Fatal(err)
	}
	return pointer
}
