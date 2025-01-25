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
	fmt.Println(world)
	// Output: world
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
	fmt.Println(string(b))
	// Output: ["hello","world"]
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

func TestValue(t *testing.T) {
	t.Run("switch syntax", func(t *testing.T) {
		target := `{"hello":"world"}`
		var value fluffyjson.Value
		if err := json.Unmarshal([]byte(target), &value); err != nil {
			t.Fatal(err)
		}

		switch object := value.Value.(type) {
		// case SomeType:
		// 	t.Fatal("fail to compile: the interface is not implemented for SomeType basically")
		case *fluffyjson.Object:
			switch world := (*object)["hello"].(type) {
			case *fluffyjson.String:
				if *world != "world" {
					t.Fatal("not world")
				}
			default:
				t.Fatal("not string")
			}
		default:
			t.Fatal("not object")
		}
	})

	t.Run("as methods", func(t *testing.T) {
		target := `{"hello":"world"}`
		var value fluffyjson.Value
		if err := json.Unmarshal([]byte(target), &value); err != nil {
			t.Fatal(err)
		}

		object, err := value.Value.AsObject()
		if err != nil {
			t.Fatal(err)
		}

		world, err := object["hello"].AsString()
		if err != nil {
			t.Fatal(err)
		}

		if world != "world" {
			t.Fatal("not world")
		}
	})
}
