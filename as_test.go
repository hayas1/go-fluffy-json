package fluffyjson_test

import (
	"testing"

	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func TestAsValue(t *testing.T) {
	t.Run("as object", func(t *testing.T) {
		testcases := map[string]struct {
			target string
			expect fluffyjson.Object
			err    error
		}{
			"object as object": {
				target: `{"hello":"world"}`,
				expect: fluffyjson.Object{"hello": fluffyjson.ForceString("world")},
				err:    nil,
			},
			"string as object": {
				target: `[0,1,2,{"three":4}]`,
				expect: nil,
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.OBJECT, Actual: fluffyjson.ARRAY},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsObject()
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as array", func(t *testing.T) {
		testcases := map[string]struct {
			target string
			expect fluffyjson.Array
			err    error
		}{
			"array as array": {
				target: `["hello", "world"]`,
				expect: fluffyjson.Array{fluffyjson.ForceString("hello"), fluffyjson.ForceString("world")},
				err:    nil,
			},
			"string as array": {
				target: `"hello world"`,
				expect: nil,
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.ARRAY, Actual: fluffyjson.STRING},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsArray()
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as string", func(t *testing.T) {
		testcases := map[string]struct {
			target string
			expect fluffyjson.String
			err    error
		}{
			"string as string": {
				target: `"hello world"`,
				expect: "hello world",
				err:    nil,
			},
			"number as string": {
				target: `100`,
				expect: "",
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.STRING, Actual: fluffyjson.NUMBER},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsString()
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as number", func(t *testing.T) {
		testcases := map[string]struct {
			target string
			expect fluffyjson.Number
			err    error
		}{
			"number as number": {
				target: `100`,
				expect: 100,
				err:    nil,
			},
			"string as number": {
				target: `true`,
				expect: 0,
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.NUMBER, Actual: fluffyjson.BOOL},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsNumber()
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as bool", func(t *testing.T) {
		testcases := map[string]struct {
			target string
			expect fluffyjson.Bool
			err    error
		}{
			"bool as bool": {
				target: `true`,
				expect: true,
				err:    nil,
			},
			"null as bool": {
				target: `null`,
				expect: false,
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.BOOL, Actual: fluffyjson.NULL},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsBool()
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as null", func(t *testing.T) {
		testcases := map[string]struct {
			target string
			expect fluffyjson.Null
			err    error
		}{
			"null as null": {
				target: `null`,
				expect: nil,
				err:    nil,
			},
			"object as null": {
				target: `{"hello": "world"}`,
				expect: []struct{}{}, // TODO nil ... ?
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.NULL, Actual: fluffyjson.OBJECT},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsNull()
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})
}

func TestAccessAsValue(t *testing.T) {
	t.Run("as object", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Object
			err      error
		}{
			"root": {
				target:   `{"hello":"world"}`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   fluffyjson.Object{"hello": fluffyjson.ForceString("world")},
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":4}]`,
				accessor: fluffyjson.ParsePointer("/3"),
				expect:   fluffyjson.Object{"three": fluffyjson.ForceNumber(4)},
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsObject(tc.accessor)
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as array", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Array
			err      error
		}{
			"root": {
				target:   `["hello", "world"]`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   fluffyjson.Array{fluffyjson.ForceString("hello"), fluffyjson.ForceString("world")},
				err:      nil,
			},
			"nested": {
				target:   `{"zero": 1, "two": [3,4]}`,
				accessor: fluffyjson.ParsePointer("/two"),
				expect:   fluffyjson.Array{fluffyjson.ForceNumber(3), fluffyjson.ForceNumber(4)},
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsArray(tc.accessor)
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as string", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.String
			err      error
		}{
			"root": {
				target:   `"hello world"`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   *fluffyjson.ForceString("hello world"),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":"4"}]`,
				accessor: fluffyjson.ParsePointer("/3/three"),
				expect:   *fluffyjson.ForceString("4"),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsString(tc.accessor)
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as number", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Number
			err      error
		}{
			"root": {
				target:   `100`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   *fluffyjson.ForceNumber(100),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":4}]`,
				accessor: fluffyjson.ParsePointer("/2"),
				expect:   *fluffyjson.ForceNumber(2),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsNumber(tc.accessor)
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as bool", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Bool
			err      error
		}{
			"root": {
				target:   `true`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   *fluffyjson.ForceBool(true),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":false}]`,
				accessor: fluffyjson.ParsePointer("/3/three"),
				expect:   *fluffyjson.ForceBool(false),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsBool(tc.accessor)
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})

	t.Run("as null", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Null
			err      error
		}{
			"root": {
				target:   `null`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   *fluffyjson.ForceNull(nil),
				err:      nil,
			},
			"nested": {
				target:   `[null,1,2,{"three":false}]`,
				accessor: fluffyjson.ParsePointer("/0"),
				expect:   *fluffyjson.ForceNull(nil),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsNull(tc.accessor)
				HelperFatalEvaluateError(t, tc.expect, actual, tc.err, err)
			})
		}
	})
}
