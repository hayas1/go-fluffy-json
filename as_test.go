package fluffyjson_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func TestAsValue(t *testing.T) {
	t.Run("as object", func(t *testing.T) {
		testcases := []struct {
			name   string
			target string
			expect fluffyjson.Object
			err    error
		}{
			{
				name:   "object as object",
				target: `{"hello":"world"}`,
				expect: fluffyjson.Object{"hello": fluffyjson.ForceString("world")},
				err:    nil,
			},
			{
				name:   "array as object",
				target: `[0,1,2,{"three":4}]`,
				expect: nil,
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.OBJECT, Actual: fluffyjson.ARRAY},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsObject()
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as array", func(t *testing.T) {
		testcases := []struct {
			name   string
			target string
			expect fluffyjson.Array
			err    error
		}{
			{
				name:   "array as array",
				target: `["hello", "world"]`,
				expect: fluffyjson.Array{fluffyjson.ForceString("hello"), fluffyjson.ForceString("world")},
				err:    nil,
			},
			{
				name:   "string as array",
				target: `"hello world"`,
				expect: nil,
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.ARRAY, Actual: fluffyjson.STRING},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsArray()
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as string", func(t *testing.T) {
		testcases := []struct {
			name   string
			target string
			expect fluffyjson.String
			err    error
		}{
			{
				name:   "string as string",
				target: `"hello world"`,
				expect: "hello world",
				err:    nil,
			},
			{
				name:   "number as string",
				target: `100`,
				expect: "",
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.STRING, Actual: fluffyjson.NUMBER},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsString()
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as number", func(t *testing.T) {
		testcases := []struct {
			name   string
			target string
			expect fluffyjson.Number
			err    error
		}{
			{
				name:   "number as number",
				target: `100`,
				expect: 100,
				err:    nil,
			},
			{
				name:   "bool as number",
				target: `true`,
				expect: 0,
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.NUMBER, Actual: fluffyjson.BOOL},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsNumber()
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as bool", func(t *testing.T) {
		testcases := []struct {
			name   string
			target string
			expect fluffyjson.Bool
			err    error
		}{
			{
				name:   "bool as bool",
				target: `true`,
				expect: true,
				err:    nil,
			},
			{
				name:   "null as bool",
				target: `null`,
				expect: false,
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.BOOL, Actual: fluffyjson.NULL},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsBool()
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as null", func(t *testing.T) {
		testcases := []struct {
			name   string
			target string
			expect fluffyjson.Null
			err    error
		}{
			{
				name:   "null as null",
				target: `null`,
				expect: nil,
				err:    nil,
			},
			{
				name:   "object as null",
				target: `{"hello": "world"}`,
				expect: []struct{}{}, // TODO nil ... ?
				err:    fluffyjson.ErrAsValue{Expected: fluffyjson.NULL, Actual: fluffyjson.OBJECT},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsNull()
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})
}

func TestAccessAsValue(t *testing.T) {
	t.Run("as object", func(t *testing.T) {
		testcases := []struct {
			name     string
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Object
			err      error
		}{
			{
				name:     "root",
				target:   `{"hello":"world"}`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   fluffyjson.Object{"hello": fluffyjson.ForceString("world")},
				err:      nil,
			},
			{
				name:     "nested",
				target:   `[0,1,2,{"three":4}]`,
				accessor: fluffyjson.ParsePointer("/3"),
				expect:   fluffyjson.Object{"three": fluffyjson.ForceNumber(4)},
				err:      nil,
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsObject(tc.accessor)
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as array", func(t *testing.T) {
		testcases := []struct {
			name     string
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Array
			err      error
		}{
			{
				name:     "root",
				target:   `["hello", "world"]`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   fluffyjson.Array{fluffyjson.ForceString("hello"), fluffyjson.ForceString("world")},
				err:      nil,
			},
			{
				name:     "nested",
				target:   `{"zero": 1, "two": [3,4]}`,
				accessor: fluffyjson.ParsePointer("/two"),
				expect:   fluffyjson.Array{fluffyjson.ForceNumber(3), fluffyjson.ForceNumber(4)},
				err:      nil,
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsArray(tc.accessor)
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as string", func(t *testing.T) {
		testcases := []struct {
			name     string
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.String
			err      error
		}{
			{
				name:     "root",
				target:   `"hello world"`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   *fluffyjson.ForceString("hello world"),
				err:      nil,
			},
			{
				name:     "nested",
				target:   `[0,1,2,{"three":"4"}]`,
				accessor: fluffyjson.ParsePointer("/3/three"),
				expect:   *fluffyjson.ForceString("4"),
				err:      nil,
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsString(tc.accessor)
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as number", func(t *testing.T) {
		testcases := []struct {
			name     string
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Number
			err      error
		}{
			{
				name:     "root",
				target:   `100`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   *fluffyjson.ForceNumber(100),
				err:      nil,
			},
			{
				name:     "nested",
				target:   `[0,1,2,{"three":4}]`,
				accessor: fluffyjson.ParsePointer("/2"),
				expect:   *fluffyjson.ForceNumber(2),
				err:      nil,
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsNumber(tc.accessor)
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as bool", func(t *testing.T) {
		testcases := []struct {
			name     string
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Bool
			err      error
		}{
			{
				name:     "root",
				target:   `true`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   *fluffyjson.ForceBool(true),
				err:      nil,
			},
			{
				name:     "nested",
				target:   `[0,1,2,{"three":false}]`,
				accessor: fluffyjson.ParsePointer("/3/three"),
				expect:   *fluffyjson.ForceBool(false),
				err:      nil,
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsBool(tc.accessor)
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("as null", func(t *testing.T) {
		testcases := []struct {
			name     string
			target   string
			accessor fluffyjson.Accessor
			expect   fluffyjson.Null
			err      error
		}{
			{
				name:     "root",
				target:   `null`,
				accessor: fluffyjson.ParsePointer("/"),
				expect:   *fluffyjson.ForceNull(nil),
				err:      nil,
			},
			{
				name:     "nested",
				target:   `[null,1,2,{"three":false}]`,
				accessor: fluffyjson.ParsePointer("/0"),
				expect:   *fluffyjson.ForceNull(nil),
				err:      nil,
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				actual, err := value.AccessAsNull(tc.accessor)
				if !errors.Is(err, tc.err) {
					t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
				} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})
}
