package fluffyjson_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	fluffyjson "github.com/hayas1/go-fluffy-json"
)

func ExampleRootValue_UnmarshalJSON() {
	var v struct {
		Fluffy fluffyjson.RootValue `json:"fluffy"`
	}
	if err := json.Unmarshal([]byte(`{"fluffy":{"deep":{"nested":{"json":{"value":["hello","world"]}}}}}`), &v); err != nil {
		panic(err)
	}

	world, err := v.Fluffy.AccessAsString(fluffyjson.ParsePointer("/deep/nested/json/value/1"))
	if err != nil {
		panic(err)
	}
	fmt.Println(world) // Output: world
}

func ExampleRootValue_MarshalJSON() {
	v := fluffyjson.RootValue{&fluffyjson.Array{
		fluffyjson.ForceString("hello"),
		fluffyjson.ForceString("world"),
	}}
	b, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b)) // Output: ["hello","world"]
}

func ExampleRootValue_asMethods() {
	target := `{"hello":"world"}`
	var value fluffyjson.RootValue
	if err := json.Unmarshal([]byte(target), &value); err != nil {
		panic(err)
	}

	object, err := value.AccessAsObject()
	if err != nil {
		panic(err)
	}

	world, err := object["hello"].AccessAsString()
	if err != nil {
		panic(err)
	}

	match := world == "world"
	if !match {
		panic("not world")
	}

	fmt.Println(match) // Output: true
}
func ExampleRootValue_switchSyntax() {
	target := `{"hello":"world"}`
	var value fluffyjson.RootValue
	if err := json.Unmarshal([]byte(target), &value); err != nil {
		panic(err)
	}

	switch object := value.JsonValue.(type) {
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

func ExampleRootValue_DepthFirst() {
	target := `[[[1,2],[3,4]],[[5,6],[7,8]]]`
	var value fluffyjson.RootValue
	if err := json.Unmarshal([]byte(target), &value); err != nil {
		panic(err)
	}

	var sum func(v fluffyjson.JsonValue) int
	sum = func(v fluffyjson.JsonValue) int {
		switch t := v.(type) {
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
	Fluffy fluffyjson.RootValue `json:"fluffy"`
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
				Fluffy: fluffyjson.RootValue{&fluffyjson.Object{"hoge": fluffyjson.ForceString("fuga")}},
			},
			err: nil,
		},
		{
			name:   "compound",
			actual: TestFluffy{},
			target: `{"fluffy":[null, true, {"three": 4}, "five"]}`,
			expect: TestFluffy{
				Fluffy: fluffyjson.RootValue{&fluffyjson.Array{
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
			HelperEvaluateError(t, tc.expect, tc.actual, tc.err, err)
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
				Fluffy: fluffyjson.RootValue{&fluffyjson.Object{"hoge": fluffyjson.ForceString("fuga")}},
			},
			expect: `{"fluffy":{"hoge":"fuga"}}`,
			err:    nil,
		},
		{
			name: "compound",
			actual: TestFluffy{
				Fluffy: fluffyjson.RootValue{&fluffyjson.Array{
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
			HelperEvaluateError(t, tc.expect, string(bytes), tc.err, err)
		})
	}
}

func HelperUnmarshalValue(t *testing.T, target string) fluffyjson.RootValue {
	t.Helper()
	var value fluffyjson.RootValue
	if err := json.Unmarshal([]byte(target), &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func HelperEvaluate[V any](t *testing.T, exp, act V) {
	t.Helper()
	HelperEvaluateError(t, exp, act, nil, nil)
}
func HelperEvaluateError[V any](t *testing.T, exp, act V, expErr, actErr error) {
	t.Helper()
	if !errors.Is(expErr, actErr) {
		t.Fatal(fmt.Errorf("%w <-> %w", expErr, actErr))
	} else if diff := cmp.Diff(exp, act); diff != "" {
		t.Fatal(diff)
	}
}
