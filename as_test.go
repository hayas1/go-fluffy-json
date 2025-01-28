package fluffyjson_test

import (
	"testing"

	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func TestAsValue(t *testing.T) {
	t.Run("as object", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			expected fluffyjson.Object
			err      error
		}{
			"object as object": {
				target:   `{"hello":"world"}`,
				expected: fluffyjson.Object{"hello": HelperCastString(t, "world")},
				err:      nil,
			},
			"string as object": {
				target:   `[0,1,2,{"three":4}]`,
				expected: nil,
				err:      fluffyjson.ErrAsValue{Expected: fluffyjson.OBJECT, Actual: fluffyjson.ARRAY},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsObject()
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as array", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			expected fluffyjson.Array
			err      error
		}{
			"array as array": {
				target:   `["hello", "world"]`,
				expected: fluffyjson.Array{HelperCastString(t, "hello"), HelperCastString(t, "world")},
				err:      nil,
			},
			"string as array": {
				target:   `"hello world"`,
				expected: nil,
				err:      fluffyjson.ErrAsValue{Expected: fluffyjson.ARRAY, Actual: fluffyjson.STRING},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsArray()
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as string", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			expected fluffyjson.String
			err      error
		}{
			"string as string": {
				target:   `"hello world"`,
				expected: "hello world",
				err:      nil,
			},
			"number as string": {
				target:   `100`,
				expected: "",
				err:      fluffyjson.ErrAsValue{Expected: fluffyjson.STRING, Actual: fluffyjson.NUMBER},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsString()
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as number", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			expected fluffyjson.Number
			err      error
		}{
			"number as number": {
				target:   `100`,
				expected: 100,
				err:      nil,
			},
			"string as number": {
				target:   `true`,
				expected: 0,
				err:      fluffyjson.ErrAsValue{Expected: fluffyjson.NUMBER, Actual: fluffyjson.BOOL},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsNumber()
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as bool", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			expected fluffyjson.Bool
			err      error
		}{
			"bool as bool": {
				target:   `true`,
				expected: true,
				err:      nil,
			},
			"null as bool": {
				target:   `null`,
				expected: false,
				err:      fluffyjson.ErrAsValue{Expected: fluffyjson.BOOL, Actual: fluffyjson.NULL},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsBool()
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as null", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			expected fluffyjson.Null
			err      error
		}{
			"null as null": {
				target:   `null`,
				expected: nil,
				err:      nil,
			},
			"object as null": {
				target:   `{"hello": "world"}`,
				expected: nil,
				err:      fluffyjson.ErrAsValue{Expected: fluffyjson.NULL, Actual: fluffyjson.OBJECT},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsNull()
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})
}

func TestAccessAsValue(t *testing.T) {
	t.Run("as object", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expected fluffyjson.Object
			err      error
		}{
			"root": {
				target:   `{"hello":"world"}`,
				accessor: fluffyjson.ParsePointer("/"),
				expected: fluffyjson.Object{"hello": HelperCastString(t, "world")},
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":4}]`,
				accessor: fluffyjson.ParsePointer("/3"),
				expected: fluffyjson.Object{"three": HelperCastNumber(t, 4)},
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsObject(tc.accessor)
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as array", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expected fluffyjson.Array
			err      error
		}{
			"root": {
				target:   `["hello", "world"]`,
				accessor: fluffyjson.ParsePointer("/"),
				expected: fluffyjson.Array{HelperCastString(t, "hello"), HelperCastString(t, "world")},
				err:      nil,
			},
			"nested": {
				target:   `{"zero": 1, "two": [3,4]}`,
				accessor: fluffyjson.ParsePointer("/two"),
				expected: fluffyjson.Array{HelperCastNumber(t, 3), HelperCastNumber(t, 4)},
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsArray(tc.accessor)
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as string", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expected fluffyjson.String
			err      error
		}{
			"root": {
				target:   `"hello world"`,
				accessor: fluffyjson.ParsePointer("/"),
				expected: *HelperCastString(t, "hello world"),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":"4"}]`,
				accessor: fluffyjson.ParsePointer("/3/three"),
				expected: *HelperCastString(t, "4"),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsString(tc.accessor)
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as number", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expected fluffyjson.Number
			err      error
		}{
			"root": {
				target:   `100`,
				accessor: fluffyjson.ParsePointer("/"),
				expected: *HelperCastNumber(t, 100),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":4}]`,
				accessor: fluffyjson.ParsePointer("/2"),
				expected: *HelperCastNumber(t, 2),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsNumber(tc.accessor)
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as bool", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expected fluffyjson.Bool
			err      error
		}{
			"root": {
				target:   `true`,
				accessor: fluffyjson.ParsePointer("/"),
				expected: *HelperCastBool(t, true),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":false}]`,
				accessor: fluffyjson.ParsePointer("/3/three"),
				expected: *HelperCastBool(t, false),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsBool(tc.accessor)
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})

	t.Run("as null", func(t *testing.T) {
		testcases := map[string]struct {
			target   string
			accessor fluffyjson.Accessor
			expected fluffyjson.Null
			err      error
		}{
			"root": {
				target:   `null`,
				accessor: fluffyjson.ParsePointer("/"),
				expected: *HelperCastNull(t, nil),
				err:      nil,
			},
			"nested": {
				target:   `[null,1,2,{"three":false}]`,
				accessor: fluffyjson.ParsePointer("/0"),
				expected: *HelperCastNull(t, nil),
				err:      nil,
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)
				actual, err := value.AccessAsNull(tc.accessor)
				HelperFatalEvaluateError(t, tc.expected, actual, tc.err, err)
			})
		}
	})
}
