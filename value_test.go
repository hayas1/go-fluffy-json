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
		var v Value
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			t.Fatal(err)
		}

		switch o := v.Value.(type) {
		// case SomeType:
		// 	t.Fatal("fail to compile: the interface is not implemented for SomeType basically")
		case *Object:
			switch world := (*o)["hello"].(type) {
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
}
