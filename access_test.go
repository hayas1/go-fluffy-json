package fluffyjson_test

import (
	"encoding/json"
	"testing"

	fluffyjson "github.com/hayas1/go-fluffy-json"
)

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
				accessor: HelperFatalParsePointer(t, "/"),
				expected: fluffyjson.Object{"hello": HelperCastString(t, "world")},
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":4}]`,
				accessor: HelperFatalParsePointer(t, "/3"),
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
				accessor: HelperFatalParsePointer(t, "/"),
				expected: fluffyjson.Array{HelperCastString(t, "hello"), HelperCastString(t, "world")},
				err:      nil,
			},
			"nested": {
				target:   `{"zero": 1, "two": [3,4]}`,
				accessor: HelperFatalParsePointer(t, "/two"),
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
				accessor: HelperFatalParsePointer(t, "/"),
				expected: *HelperCastString(t, "hello world"),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":"4"}]`,
				accessor: HelperFatalParsePointer(t, "/3/three"),
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
				accessor: HelperFatalParsePointer(t, "/"),
				expected: *HelperCastNumber(t, 100),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":4}]`,
				accessor: HelperFatalParsePointer(t, "/2"),
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
				accessor: HelperFatalParsePointer(t, "/"),
				expected: *HelperCastBool(t, true),
				err:      nil,
			},
			"nested": {
				target:   `[0,1,2,{"three":false}]`,
				accessor: HelperFatalParsePointer(t, "/3/three"),
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
				accessor: HelperFatalParsePointer(t, "/"),
				expected: *HelperCastNull(t, nil),
				err:      nil,
			},
			"nested": {
				target:   `[null,1,2,{"three":false}]`,
				accessor: HelperFatalParsePointer(t, "/0"),
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

func TestSubtreeReplace(t *testing.T) {
	t.Run("subtree replace", func(t *testing.T) {
		value := HelperUnmarshalValue(t, `[0,1,2]`)
		array, err := value.AccessAsArray()
		if err != nil {
			t.Fatal(err)
		}
		hello := HelperUnmarshalValue(t, `{"hello":"world"}`)
		array[1] = &hello

		replaced, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		expected := `[0,{"hello":"world"},2]`
		HelperFatalEvaluate(t, expected, string(replaced))
	})
}

func FuzzAccessAsValueAndAsValue(f *testing.F) {
	f.Add(`{"hello":"world"}`)
	f.Add(`[1,2,3]`)
	f.Add(`"string"`)
	f.Add(`100`)
	f.Add(`true`)
	f.Add(`null`)
	f.Fuzz(func(t *testing.T, target string) {
		var value fluffyjson.RootValue
		if err := json.Unmarshal([]byte(target), &value); err != nil {
			return
		}
		asObject, errAs := value.AsObject()
		accessAsObject, errAccessAs := value.AccessAsObject()
		HelperFatalEvaluateError(t, asObject, accessAsObject, errAs, errAccessAs)

		asArray, errAs := value.AsArray()
		accessAsArray, errAccessAs := value.AccessAsArray()
		HelperFatalEvaluateError(t, asArray, accessAsArray, errAs, errAccessAs)

		asString, errAs := value.AsString()
		accessAsString, errAccessAs := value.AccessAsString()
		HelperFatalEvaluateError(t, asString, accessAsString, errAs, errAccessAs)

		asNumber, errAs := value.AsNumber()
		accessAsNumber, errAccessAs := value.AccessAsNumber()
		HelperFatalEvaluateError(t, asNumber, accessAsNumber, errAs, errAccessAs)

		asBool, errAs := value.AsBool()
		accessAsBool, errAccessAs := value.AccessAsBool()
		HelperFatalEvaluateError(t, asBool, accessAsBool, errAs, errAccessAs)

		asNull, errAs := value.AsNull()
		accessAsNull, errAccessAs := value.AccessAsNull()
		HelperFatalEvaluateError(t, asNull, accessAsNull, errAs, errAccessAs)
	})
}
