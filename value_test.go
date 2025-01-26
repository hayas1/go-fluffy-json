package fluffyjson_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func ExampleValue_UnmarshalJSON() {
	var v struct {
		Fluffy fluffyjson.Value `json:"fluffy"`
	}
	if err := json.Unmarshal([]byte(`{"fluffy":{"deep":{"nested":{"json":{"value":["hello","world"]}}}}}`), &v); err != nil {
		panic(err)
	}

	worldValue, err := v.Fluffy.Access(fluffyjson.ParsePointer("/deep/nested/json/value/1"))
	if err != nil {
		panic(err)
	}
	world, err := worldValue.AsString()
	if err != nil {
		panic(err)
	}
	fmt.Println(world) // Output: world
}

func ExampleValue_MarshalJSON() {
	v := fluffyjson.Value{Value: &fluffyjson.Array{
		fluffyjson.ForceString("hello"),
		fluffyjson.ForceString("world"),
	}}
	b, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b)) // Output: ["hello","world"]
}

func ExampleValue_asMethods() {
	target := `{"hello":"world"}`
	var value fluffyjson.Value
	if err := json.Unmarshal([]byte(target), &value); err != nil {
		panic(err)
	}

	object, err := value.Value.AsObject()
	if err != nil {
		panic(err)
	}

	world, err := object["hello"].AsString()
	if err != nil {
		panic(err)
	}

	match := world == "world"
	if !match {
		panic("not world")
	}

	fmt.Println(match) // Output: true
}
func ExampleValue_switchSyntax() {
	target := `{"hello":"world"}`
	var value fluffyjson.Value
	if err := json.Unmarshal([]byte(target), &value); err != nil {
		panic(err)
	}

	switch object := value.Value.(type) {
	// case SomeType:
	// 	panic("fail to compile: the interface is not implemented for SomeType basically")
	case *fluffyjson.Object:
		switch world := (*object)["hello"].(type) {
		case *fluffyjson.String:
			if match := *world == "world"; !match {
				panic("not world")
			} else {
				fmt.Println(match)
			}
		default:
			panic("not string")
		}
	default:
		panic("not object")
	}
	// Output: true
}

func ExampleValue_DepthFirst() {
	target := `[[[1,2],[3,4]],[[5,6],[7,8]]]`
	var value fluffyjson.Value
	if err := json.Unmarshal([]byte(target), &value); err != nil {
		panic(err)
	}

	var sum func(v fluffyjson.JsonValue) int
	sum = func(v fluffyjson.JsonValue) int {
		switch t := v.(type) {
		case *fluffyjson.Value:
			return sum(t.Value) // TODO
		case *fluffyjson.Array:
			s := 0
			for _, vv := range *t {
				s += sum(vv)
			}
			return s
		case *fluffyjson.Number:
			return int(*t)
		default:
			panic("not array or number")
		}
	}
	results := make([]int, 0, 15)
	for _, v := range value.DepthFirst() {
		results = append(results, sum(v))
	}
	fmt.Println(results) // Output: [36 10 3 1 2 7 3 4 26 11 5 6 15 7 8]
}

type TestFluffy struct {
	Fluffy fluffyjson.Value `json:"fluffy"`
}

func TestUnmarshalBasic(t *testing.T) {
	testcases := []struct {
		name   string
		actual TestFluffy
		target string
		expect TestFluffy
		err    error
	}{
		{
			name:   "object and string",
			actual: TestFluffy{},
			target: `{"fluffy":{"hoge":"fuga"}}`,
			expect: TestFluffy{
				Fluffy: fluffyjson.Value{Value: &fluffyjson.Object{"hoge": fluffyjson.ForceString("fuga")}},
			},
			err: nil,
		},
		{
			name:   "compound",
			actual: TestFluffy{},
			target: `{"fluffy":[null, true, {"three": 4}, "five"]}`,
			expect: TestFluffy{
				Fluffy: fluffyjson.Value{Value: &fluffyjson.Array{
					fluffyjson.ForceNull(nil),
					fluffyjson.ForceBool(true),
					&fluffyjson.Object{"three": fluffyjson.ForceNumber(4)},
					fluffyjson.ForceString("five"),
				}},
			},
			err: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tc.target), &tc.actual)
			if !errors.Is(err, tc.err) {
				t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
			} else if diff := cmp.Diff(tc.expect, tc.actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
func TestMarshalBasic(t *testing.T) {
	testcases := []struct {
		name   string
		actual TestFluffy
		expect string
		err    error
	}{
		{
			name: "object and string",
			actual: TestFluffy{
				Fluffy: fluffyjson.Value{Value: &fluffyjson.Object{"hoge": fluffyjson.ForceString("fuga")}},
			},
			expect: `{"fluffy":{"hoge":"fuga"}}`,
			err:    nil,
		},
		{
			name: "compound",
			actual: TestFluffy{
				Fluffy: fluffyjson.Value{Value: &fluffyjson.Array{
					fluffyjson.ForceNull(nil),
					fluffyjson.ForceBool(true),
					&fluffyjson.Object{"three": fluffyjson.ForceNumber(4)},
					fluffyjson.ForceString("five"),
				}},
			},
			expect: `{"fluffy":[null,true,{"three":4},"five"]}`,
			err:    nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := json.Marshal(&tc.actual)
			if !errors.Is(err, tc.err) {
				t.Fatal(fmt.Errorf("%w <-> %w", tc.err, err))
			} else if diff := cmp.Diff(tc.expect, string(bytes)); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
