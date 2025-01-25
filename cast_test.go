package fluffyjson_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayash1/go-fluffy-json"
)

func TestValueAs(t *testing.T) {
	testcases := []struct {
		name   string
		target string
		as     func(fluffyjson.JsonValue) (fluffyjson.JsonValue, error)
		expect string
		err    error
	}{
		{
			name:   "object as object",
			target: `{"number": ["zero", "one", "two"]}`,
			as:     func(jv fluffyjson.JsonValue) (fluffyjson.JsonValue, error) { o, e := jv.AsObject(); return &o, e },
			expect: `{"number": ["zero", "one", "two"]}`,
			err:    nil,
		},
		{
			name:   "array as string",
			target: `["hello", "world"]`,
			as:     func(jv fluffyjson.JsonValue) (fluffyjson.JsonValue, error) { a, e := jv.AsString(); return &a, e },
			expect: `""`,
			err:    fluffyjson.ErrAsValue{Not: "string", But: "array"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var value, expect fluffyjson.Value
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tc.expect), &expect); err != nil {
				t.Fatal(err)
			}

			actual, err := tc.as(&value)
			if !errors.Is(err, tc.err) {
				t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
			} else if diff := cmp.Diff(expect, fluffyjson.Value{Value: actual}); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
