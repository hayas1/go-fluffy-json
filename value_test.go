package fluffyjson

// TODO package fluffyjson_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type HelloWorld struct {
	Hello String `json:"hello"`
	Meta  Value  `json:"meta"`
}

func TestUnmarshalBasic(t *testing.T) {
	testCases := []struct {
		name   string
		actual HelloWorld
		input  string
		expect HelloWorld
		err    error
	}{
		{
			name:   "hello world",
			actual: HelloWorld{},
			// input:  `{"hello":"world"}`,
			input: `{"hello":"world", "meta":{"hoge":"fuga"}}`,
			expect: HelloWorld{
				Hello: "world",
				Meta:  Value{Value: &Object{"hoge": &[]String{("fuga")}[0]}},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tc.input), &tc.actual)
			if err != tc.err {
				t.Fatal(err)
			} else if diff := cmp.Diff(tc.expect, tc.actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
func TestMarshalBasic(t *testing.T) {
	testCases := []struct {
		name   string
		actual HelloWorld
		expect string
		err    error
	}{
		{
			name: "hello world",
			actual: HelloWorld{
				Hello: "world",
				Meta:  Value{Value: &Object{"hoge": &[]String{("fuga")}[0]}},
			},
			expect: `{"hello":"world","meta":{"hoge":"fuga"}}`,
			err:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := json.Marshal(&tc.actual)
			if err != tc.err {
				t.Fatal(err)
			} else if diff := cmp.Diff(tc.expect, string(bytes)); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestValue(t *testing.T) {
	t.Run("switch syntax", func(t *testing.T) {
		raw := `{"hello":"world"}`
		var value Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}

		switch object := value.Value.(type) {
		// case SomeType:
		// 	t.Fatal("fail to compile: the interface is not implemented for SomeType basically")
		case *Object:
			switch world := (*object)["hello"].(type) {
			case *String:
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
		raw := `{"hello":"world"}`
		var value Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
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

	t.Run("iterator", func(t *testing.T) {
		raw := `{"hello":"world"}`
		var value Value
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}

		for k, v := range value.Iter() {
			hello, err := k.AsObjectKey()
			if err != nil {
				t.Fatal(err)
			}
			if hello != "hello" {
				t.Fatal("not hello")
			}
			world, err := v.AsString()
			if err != nil {
				t.Fatal(err)
			}
			if world != "world" {
				t.Fatal("not world")
			}
		}
	})
}
