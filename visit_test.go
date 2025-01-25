// TODO package fluffyjson_test
package fluffyjson

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Collector struct {
	DfsVisitor
	visited []string
}

//	func (collector *Collector) VisitObject(s *Object) error {
//		return nil
//	}
//
//	func (collector *Collector) VisitArray(s *Array) error {
//		return nil
//	}
func (collector *Collector) VisitString(s *String) error {
	collector.visited = append(collector.visited, string(*s))
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
		// visitor := DfsVisitor{Inner: collector}
		if err := value.Accept(collector); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff([]string{"a", "b", "c", "d", "e", "f", "g"}, collector.visited); diff != "" {
			t.Fatal(diff)
		}
	})

}
