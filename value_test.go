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

func TestUnmarshal(t *testing.T) {
	type HelloWorld struct {
		Hello string `json:"hello"`
	}
	testCases := []struct {
		name   string
		actual HelloWorld
		input  string
		expect interface{}
		err    error
	}{
		{
			name:   "hello world",
			actual: HelloWorld{},
			input:  `{"hello":"world"}`,
			expect: HelloWorld{Hello: "world"},
			err:    nil,
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
