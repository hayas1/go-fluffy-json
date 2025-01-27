package fluffyjson_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func TestAsValue(t *testing.T) {
	testcases := []struct {
		name   string
		target string
		as     func(fluffyjson.JsonValue) (fluffyjson.JsonValue, error)
		expect fluffyjson.JsonValue
		err    error
	}{
		{
			name:   "object as object",
			target: `{"number": ["zero", "one", "two"]}`,
			as:     func(jv fluffyjson.JsonValue) (fluffyjson.JsonValue, error) { o, e := jv.AccessAsObject(); return &o, e },
			expect: &fluffyjson.Object{
				"number": &fluffyjson.Array{
					fluffyjson.ForceString("zero"), fluffyjson.ForceString("one"), fluffyjson.ForceString("two"),
				},
			},
			err: nil,
		},
		{
			name:   "array as string",
			target: `["hello", "world"]`,
			as:     func(jv fluffyjson.JsonValue) (fluffyjson.JsonValue, error) { a, e := jv.AccessAsString(); return &a, e },
			expect: fluffyjson.ForceString(""),
			err:    fluffyjson.ErrAsValue{Not: fluffyjson.STRING, But: fluffyjson.ARRAY},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.RootValue
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			actual, err := tc.as(&value)
			if !errors.Is(err, tc.err) {
				t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
			} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
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
		}

		for _, tc := range testcases {
			var value fluffyjson.RootValue
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			actual, err := value.AccessAsObject(tc.accessor)
			if !errors.Is(err, tc.err) {
				t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
			} else if diff := cmp.Diff(tc.expect, actual); diff != "" {
				t.Fatal(diff)
			}
		}
	})
}
