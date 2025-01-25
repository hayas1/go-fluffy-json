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

}
