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
		testcases := map[string]struct {
			target          string
			expectedPointer [][]fluffyjson.Pointer
			expectedValue   [][]fluffyjson.JsonValue
		}{
			"basic": {
				target: `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`,
				expectedPointer: [][]fluffyjson.Pointer{
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
				expectedValue: [][]fluffyjson.JsonValue{
					{
						&fluffyjson.Object{
							"a": &fluffyjson.Object{
								"b": &fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
								"e": &fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
							},
						},
						&fluffyjson.Object{
							"b": &fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
							"e": &fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
						},
						&fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
						HelperCastString(t, "c"),
						HelperCastString(t, "d"),
						&fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
						HelperCastString(t, "f"),
						HelperCastString(t, "g"),
					},
					{
						&fluffyjson.Object{
							"a": &fluffyjson.Object{
								"b": &fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
								"e": &fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
							},
						},
						&fluffyjson.Object{
							"b": &fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
							"e": &fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
						},
						&fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
						HelperCastString(t, "f"),
						HelperCastString(t, "g"),
						&fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
						HelperCastString(t, "c"),
						HelperCastString(t, "d"),
					},
				},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				var i int
				for p, v := range value.DepthFirst() {
					pAny, pDiff := false, make([]string, 0, len(tc.expectedPointer))
					for _, ep := range tc.expectedPointer {
						diff := cmp.Diff(ep[i], p)
						pDiff = append(pDiff, diff)
						pAny = pAny || diff == ""
					}
					if !pAny {
						t.Fatal(pDiff)
					}

					vAny, vDiff := false, make([]string, 0, len(tc.expectedValue))
					for _, ev := range tc.expectedValue {
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
		testcases := map[string]struct {
			target          string
			expectedPointer [][]fluffyjson.Pointer
			expectedValue   [][]fluffyjson.JsonValue
		}{
			"basic": {
				target: `{"a":{"b": ["c", "d"], "e": ["f", "g"]}}`,
				expectedPointer: [][]fluffyjson.Pointer{
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
				expectedValue: [][]fluffyjson.JsonValue{
					{
						&fluffyjson.Object{
							"a": &fluffyjson.Object{
								"b": &fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
								"e": &fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
							},
						},
						&fluffyjson.Object{
							"b": &fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
							"e": &fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
						},
						&fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
						&fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
						HelperCastString(t, "c"),
						HelperCastString(t, "d"),
						HelperCastString(t, "f"),
						HelperCastString(t, "g"),
					},
					{
						&fluffyjson.Object{
							"a": &fluffyjson.Object{
								"b": &fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
								"e": &fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
							},
						},
						&fluffyjson.Object{
							"b": &fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
							"e": &fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
						},
						&fluffyjson.Array{HelperCastString(t, "f"), HelperCastString(t, "g")},
						&fluffyjson.Array{HelperCastString(t, "c"), HelperCastString(t, "d")},
						HelperCastString(t, "f"),
						HelperCastString(t, "g"),
						HelperCastString(t, "c"),
						HelperCastString(t, "d"),
					},
				},
			},
		}

		for name, tc := range testcases {
			t.Run(name, func(t *testing.T) {
				value := HelperUnmarshalValue(t, tc.target)

				var i int
				for p, v := range value.BreadthFirst() {
					pAny, pDiff := false, make([]string, 0, len(tc.expectedPointer))
					for _, ep := range tc.expectedPointer {
						diff := cmp.Diff(ep[i], p)
						pDiff = append(pDiff, diff)
						pAny = pAny || diff == ""
					}
					if !pAny {
						t.Fatal(pDiff)
					}

					vAny, vDiff := false, make([]string, 0, len(tc.expectedValue))
					for _, ev := range tc.expectedValue {
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
