package fluffyjson_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func TestAccessAs(t *testing.T) {
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
