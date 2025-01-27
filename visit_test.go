package fluffyjson_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayas1/go-fluffy-json"
)

type Collector struct {
	fluffyjson.PointerVisitor
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
		target := `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`
		value := HelperUnmarshalValue(t, target)

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

func TestSearch(t *testing.T) {
	t.Run("depth first", func(t *testing.T) {
		testcases := []struct {
			name          string
			target        string
			expectPointer [][]fluffyjson.Pointer
			expectValue   [][]fluffyjson.JsonValue
		}{
			{
				name:   "basic",
				target: `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`,
				expectPointer: [][]fluffyjson.Pointer{
					{
						nil,
						fluffyjson.Pointer{fluffyjson.KeyAccess("a")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b"), fluffyjson.IndexAccess(0)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b"), fluffyjson.IndexAccess(1)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e"), fluffyjson.IndexAccess(0)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e"), fluffyjson.IndexAccess(1)},
					},
					{
						nil,
						fluffyjson.Pointer{fluffyjson.KeyAccess("a")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e"), fluffyjson.IndexAccess(0)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e"), fluffyjson.IndexAccess(1)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b"), fluffyjson.IndexAccess(0)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b"), fluffyjson.IndexAccess(1)},
					},
				},
				expectValue: [][]fluffyjson.JsonValue{
					{
						&fluffyjson.Object{
							"a": &fluffyjson.Object{
								"b": &fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
								"e": &fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
							},
						},
						&fluffyjson.Object{
							"b": &fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
							"e": &fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
						},
						&fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
						fluffyjson.ForceString("c"),
						fluffyjson.ForceString("d"),
						&fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
						fluffyjson.ForceString("f"),
						fluffyjson.ForceString("g"),
					},
					{
						&fluffyjson.Object{
							"a": &fluffyjson.Object{
								"b": &fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
								"e": &fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
							},
						},
						&fluffyjson.Object{
							"b": &fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
							"e": &fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
						},
						&fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
						fluffyjson.ForceString("f"),
						fluffyjson.ForceString("g"),
						&fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
						fluffyjson.ForceString("c"),
						fluffyjson.ForceString("d"),
					},
				},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				var i int
				for p, v := range value.DepthFirst() {
					pAny, pDiff := false, make([]string, 0, len(tc.expectPointer))
					for _, ep := range tc.expectPointer {
						diff := cmp.Diff(ep[i], p)
						pDiff = append(pDiff, diff)
						pAny = pAny || diff == ""
					}
					if !pAny {
						t.Fatal(pDiff)
					}

					vAny, vDiff := false, make([]string, 0, len(tc.expectValue))
					for _, ev := range tc.expectValue {
						diff := cmp.Diff(ev[i], v)
						vDiff = append(vDiff, diff)
						vAny = vAny || diff == ""
					}
					if !vAny {
						t.Fatal(vDiff)
					}
					i++
				}
			})
		}
	})

	t.Run("breadth first", func(t *testing.T) {
		testcases := []struct {
			name          string
			target        string
			expectPointer [][]fluffyjson.Pointer
			expectValue   [][]fluffyjson.JsonValue
		}{
			{
				name:   "basic",
				target: `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`,
				expectPointer: [][]fluffyjson.Pointer{
					{
						nil,
						fluffyjson.Pointer{fluffyjson.KeyAccess("a")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b"), fluffyjson.IndexAccess(0)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b"), fluffyjson.IndexAccess(1)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e"), fluffyjson.IndexAccess(0)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e"), fluffyjson.IndexAccess(1)},
					},
					{
						nil,
						fluffyjson.Pointer{fluffyjson.KeyAccess("a")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b")},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e"), fluffyjson.IndexAccess(0)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("e"), fluffyjson.IndexAccess(1)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b"), fluffyjson.IndexAccess(0)},
						fluffyjson.Pointer{fluffyjson.KeyAccess("a"), fluffyjson.KeyAccess("b"), fluffyjson.IndexAccess(1)},
					},
				},
				expectValue: [][]fluffyjson.JsonValue{
					{
						&fluffyjson.Object{
							"a": &fluffyjson.Object{
								"b": &fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
								"e": &fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
							},
						},
						&fluffyjson.Object{
							"b": &fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
							"e": &fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
						},
						&fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
						&fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
						fluffyjson.ForceString("c"),
						fluffyjson.ForceString("d"),
						fluffyjson.ForceString("f"),
						fluffyjson.ForceString("g"),
					},
					{
						&fluffyjson.Object{
							"a": &fluffyjson.Object{
								"b": &fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
								"e": &fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
							},
						},
						&fluffyjson.Object{
							"b": &fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
							"e": &fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
						},
						&fluffyjson.Array{fluffyjson.ForceString("f"), fluffyjson.ForceString("g")},
						&fluffyjson.Array{fluffyjson.ForceString("c"), fluffyjson.ForceString("d")},
						fluffyjson.ForceString("f"),
						fluffyjson.ForceString("g"),
						fluffyjson.ForceString("c"),
						fluffyjson.ForceString("d"),
					},
				},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				var i int
				for p, v := range value.BreadthFirst() {
					pAny, pDiff := false, make([]string, 0, len(tc.expectPointer))
					for _, ep := range tc.expectPointer {
						diff := cmp.Diff(ep[i], p)
						pDiff = append(pDiff, diff)
						pAny = pAny || diff == ""
					}
					if !pAny {
						t.Fatal(pDiff)
					}

					vAny, vDiff := false, make([]string, 0, len(tc.expectValue))
					for _, ev := range tc.expectValue {
						diff := cmp.Diff(ev[i], v)
						vDiff = append(vDiff, diff)
						vAny = vAny || diff == ""
					}
					if !vAny {
						t.Fatal(vDiff)
					}
					i++
				}
			})
		}
	})
}
