package fluffyjson_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayash1/go-fluffy-json"
)

type Collector struct {
	fluffyjson.BaseVisitor
	visited string
}

func (c *Collector) VisitObjectEntry(k string, v fluffyjson.JsonValue) error {
	c.visited += k
	return nil
}
func (c *Collector) VisitString(s *fluffyjson.String) error {
	c.visited += string(*s)
	return nil
}

func TestDfsVisitor(t *testing.T) {
	t.Run("dfs collect visitor", func(t *testing.T) {
		raw := `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`
		var value fluffyjson.Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}

		collector := &Collector{}
		if err := value.Accept(fluffyjson.DfsVisitor(collector)); err != nil {
			t.Fatal(err)
		}

		diff1, diff2 := cmp.Diff("abcdefg", collector.visited), cmp.Diff("aefgbcd", collector.visited)
		if (diff1 != "" || diff2 == "") && (diff1 == "" || diff2 != "") {
			t.Fatal(diff1, diff2)
		}
	})
}

func TestDepthFirst(t *testing.T) {
	testCases := []struct {
		name          string
		target        string
		expectPointer []fluffyjson.Pointer
		expect        []fluffyjson.JsonValue
		err           error
	}{
		{
			name:   "depth first",
			target: `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`,
			expectPointer: []fluffyjson.Pointer{
				fluffyjson.ParsePointer("/a"),
				fluffyjson.ParsePointer("/a/b"),
				fluffyjson.ParsePointer("/a/b/0"),
				fluffyjson.ParsePointer("/a/b/1"),
				fluffyjson.ParsePointer("/a/e"),
				fluffyjson.ParsePointer("/a/e/0"),
				fluffyjson.ParsePointer("/a/e/1"),
			},
			expect: []fluffyjson.JsonValue{
				&fluffyjson.Object{
					"b": &fluffyjson.Array{&[]fluffyjson.String{("c")}[0], &[]fluffyjson.String{("d")}[0]},
					"e": &fluffyjson.Array{&[]fluffyjson.String{("f")}[0], &[]fluffyjson.String{("g")}[0]},
				},
				&fluffyjson.Array{&[]fluffyjson.String{("c")}[0], &[]fluffyjson.String{("d")}[0]},
				&[]fluffyjson.String{("c")}[0],
				&[]fluffyjson.String{("d")}[0],
				&fluffyjson.Array{&[]fluffyjson.String{("f")}[0], &[]fluffyjson.String{("g")}[0]},
				&[]fluffyjson.String{("f")}[0],
				&[]fluffyjson.String{("g")}[0],
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var value fluffyjson.Value
			if err := json.Unmarshal([]byte(tc.target), &value); err != nil {
				t.Fatal(err)
			}

			var i int
			for p, v := range value.DepthFirst() {
				if diff := cmp.Diff(tc.expectPointer[i], p); diff != "" {
					t.Fatal(diff)
				}
				if diff := cmp.Diff(tc.expect[i], v); diff != "" {
					t.Fatal(diff)
				}
				i++
			}
		})
	}
}
