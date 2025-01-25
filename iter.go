package fluffyjson

import (
	"fmt"
	"iter"
)

type Shallow interface {
	Iter() iter.Seq2[string, JsonValue]
}

func (v Value) Iter() iter.Seq2[string, JsonValue] {
	return v.Value.Iter()
}

func (o Object) Iter() iter.Seq2[string, JsonValue] {
	return func(yield func(string, JsonValue) bool) {
		for k, v := range o {
			yield(k, v)
		}
	}
}

func (a Array) Iter() iter.Seq2[string, JsonValue] {
	return func(yield func(string, JsonValue) bool) {
		for i, v := range a {
			yield(fmt.Sprint(i), v)
		}
	}
}

func (s String) Iter() iter.Seq2[string, JsonValue] {
	return func(yield func(string, JsonValue) bool) {
		yield("", &s)
	}
}
