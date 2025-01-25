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
	t.Run("dfs visitor", func(t *testing.T) {
		raw := `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`
		var value fluffyjson.Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}

		collector := &Collector{}
		if err := value.Accept(fluffyjson.DfsVisitor(collector)); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff("abcdefg", collector.visited); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("dfs iterator", func(t *testing.T) {
		raw := `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`
		var value fluffyjson.Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}

		expect_path := [][]string{
			{"a"}, {"a", "b"}, {"a", "b", "0"}, {"a", "b", "1"},
			{"a", "e"}, {"a", "e", "0"}, {"a", "e", "1"},
		}
		expect := []fluffyjson.JsonValue{
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
		}
		var i int
		for p, v := range value.DepthFirst() {
			if diff := cmp.Diff(expect_path[i], p); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(expect[i], v); diff != "" {
				t.Fatal(diff)
			}
			i++
		}
	})
}
