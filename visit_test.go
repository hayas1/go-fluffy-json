// TODO package fluffyjson_test
package fluffyjson

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Collector struct {
	BaseVisitor
	visited string
}

func (c *Collector) VisitObjectEntry(k string, v JsonValue) error {
	c.visited += k
	return nil
}
func (c *Collector) VisitString(s *String) error {
	c.visited += string(*s)
	return nil
}

func TestDfsVisitor(t *testing.T) {
	t.Run("dfs visitor", func(t *testing.T) {
		raw := `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`
		var value Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}

		collector := &Collector{}
		if err := value.Accept(DfsVisitor(collector)); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff("abcdefg", collector.visited); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("dfs iterator", func(t *testing.T) {
		raw := `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`
		var value Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}

		expect_path := [][]string{
			{"a"}, {"a", "b"}, {"a", "b", "0"}, {"a", "b", "1"},
			{"a", "e"}, {"a", "e", "0"}, {"a", "e", "1"},
		}
		expect := []JsonValue{
			&Object{
				"b": &Array{&[]String{("c")}[0], &[]String{("d")}[0]},
				"e": &Array{&[]String{("f")}[0], &[]String{("g")}[0]},
			},
			&Array{&[]String{("c")}[0], &[]String{("d")}[0]},
			&[]String{("c")}[0],
			&[]String{("d")}[0],
			&Array{&[]String{("f")}[0], &[]String{("g")}[0]},
			&[]String{("f")}[0],
			&[]String{("g")}[0],
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
