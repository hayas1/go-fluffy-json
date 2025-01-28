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

	pointer, err := fluffyjson.ParsePointer("/deep/nested/json/value/1")
	if err != nil {
		panic(err)
	}

	world, err := v.Fluffy.AccessAsString(pointer)
	if err != nil {
		panic(err)
	}
	fmt.Println(world) // Output: world
}

func ExampleRootValue_MarshalJSON() {
	v := fluffyjson.RootValue{&fluffyjson.Array{
		&[]fluffyjson.String{fluffyjson.String("hello")}[0],
		&[]fluffyjson.String{fluffyjson.String("world")}[0],
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
	testcases := map[string]struct {
		actual   TestFluffy
		target   string
		expected TestFluffy
		err      error
	}{
		"object and string": {
			actual: TestFluffy{},
			target: `{"fluffy":{"hoge":"fuga"}}`,
			expected: TestFluffy{
				Fluffy: fluffyjson.RootValue{&fluffyjson.Object{"hoge": HelperCastString(t, "fuga")}},
			},
			err: nil,
		},
		"compound": {
			actual: TestFluffy{},
			target: `{"fluffy":[null, true, {"three": 4}, "five"]}`,
			expected: TestFluffy{
				Fluffy: fluffyjson.RootValue{&fluffyjson.Array{
					HelperCastNull(t, nil),
					HelperCastBool(t, true),
					&fluffyjson.Object{"three": HelperCastNumber(t, 4)},
					HelperCastString(t, "five"),
				}},
			},
			err: nil,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tc.target), &tc.actual)
			HelperFatalEvaluateError(t, tc.expected, tc.actual, tc.err, err)
		})
	}
}
func TestMarshalBasic(t *testing.T) {
	testcases := map[string]struct {
		actual   TestFluffy
		expected string
		err      error
	}{
		"object and string": {
			actual: TestFluffy{
				Fluffy: fluffyjson.RootValue{&fluffyjson.Object{"hoge": HelperCastString(t, "fuga")}},
			},
			expected: `{"fluffy":{"hoge":"fuga"}}`,
			err:      nil,
		},
		"compound": {
			actual: TestFluffy{
				Fluffy: fluffyjson.RootValue{&fluffyjson.Array{
					HelperCastNull(t, nil),
					HelperCastBool(t, true),
					&fluffyjson.Object{"three": HelperCastNumber(t, 4)},
					HelperCastString(t, "five"),
				}},
			},
			expected: `{"fluffy":[null,true,{"three":4},"five"]}`,
			err:      nil,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			bytes, err := json.Marshal(&tc.actual)
			HelperFatalEvaluateError(t, tc.expected, string(bytes), tc.err, err)
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
func HelperMarshalValue(t *testing.T, value fluffyjson.RootValue) string {
	t.Helper()
	bytes, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

func HelperFatalEvaluate[V any](t *testing.T, exp, act V) {
	t.Helper()
	HelperFatalEvaluateError(t, exp, act, nil, nil)
}
func HelperFatalEvaluateError[V any](t *testing.T, exp, act V, expErr, actErr error) {
	t.Helper()
	if fatal := HelperEvaluateError(t, exp, act, expErr, actErr); fatal != nil {
		t.Fatal(fatal...)
	}
}
func HelperEvaluateError[V any](t *testing.T, exp, act V, expErr, actErr error) []any {
	t.Helper()
	if !errors.Is(expErr, actErr) {
		return []any{expErr, actErr}
	}
	if diff := cmp.Diff(exp, act); diff != "" {
		return []any{diff}
	}
	return nil
}
func HelperCastObject(t *testing.T, m map[string]any) fluffyjson.Object {
	t.Helper()
	v, err := fluffyjson.CastObject(m)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func HelperCastArray(t *testing.T, l []any) fluffyjson.Array {
	t.Helper()
	v, err := fluffyjson.CastArray(l)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func HelperCastString(t *testing.T, s string) *fluffyjson.String {
	t.Helper()
	v, err := fluffyjson.CastString(s)
	if err != nil {
		t.Fatal(err)
	}
	return &v
}
func HelperCastNumber(t *testing.T, n float64) *fluffyjson.Number {
	t.Helper()
	v, err := fluffyjson.CastNumber(n)
	if err != nil {
		t.Fatal(err)
	}
	return &v
}
func HelperCastBool(t *testing.T, b bool) *fluffyjson.Bool {
	t.Helper()
	v, err := fluffyjson.CastBool(b)
	if err != nil {
		t.Fatal(err)
	}
	return &v
}
func HelperCastNull(t *testing.T, n func()) *fluffyjson.Null {
	// func HelperCastNull(t *testing.T, n func(null)) *fluffyjson.Null {
	t.Helper()
	if n != nil {
		t.Fatalf("%p is not nil", n)
	}
	v, err := fluffyjson.CastNull(nil)
	if err != nil {
		t.Fatal(err)
	}
	return &v
}
