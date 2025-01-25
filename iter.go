package fluffyjson

import (
	"fmt"
	"iter"
)

type (
	Shallow interface {
		Iter() iter.Seq2[IterKey, JsonValue]
	}
)

func (v Value) Iter() iter.Seq2[IterKey, JsonValue] {
	return v.Value.Iter()
}

func (o Object) Iter() iter.Seq2[IterKey, JsonValue] {
	return func(yield func(IterKey, JsonValue) bool) {
		for k, v := range o {
			if !yield(ObjectKey(k), v) {
				break
			}
		}
	}
}

func (a Array) Iter() iter.Seq2[IterKey, JsonValue] {
	return func(yield func(IterKey, JsonValue) bool) {
		for i, v := range a {
			if !yield(ArrayIndex(i), v) {
				break
			}
		}
	}
}

func (s String) Iter() iter.Seq2[IterKey, JsonValue] {
	return func(yield func(IterKey, JsonValue) bool) {
		yield(Unit{}, &s)
	}
}

type (
	IterKey interface {
		IsObjectKey() bool
		AsObjectKey() (ObjectKey, error)
		IsArrayIndex() bool
		AsArrayIndex() (ArrayIndex, error)
	}

	ObjectKey  string
	ArrayIndex int
	Unit       struct{}
)

func (ok ObjectKey) AsObjectKey() (ObjectKey, error) { return ok, nil }
func (ok ObjectKey) IsObjectKey() bool               { return true }
func (ok ObjectKey) IsArrayIndex() bool              { return false }
func (ok ObjectKey) AsArrayIndex() (ArrayIndex, error) {
	return 0, fmt.Errorf("not array index, but object key")
}

func (ai ArrayIndex) IsObjectKey() bool { return false }
func (ai ArrayIndex) AsObjectKey() (ObjectKey, error) {
	return "", fmt.Errorf("not object key, but array index")
}
func (ai ArrayIndex) IsArrayIndex() bool                { return true }
func (ai ArrayIndex) AsArrayIndex() (ArrayIndex, error) { return ai, nil }

func (u Unit) IsObjectKey() bool { return false }
func (u Unit) AsObjectKey() (ObjectKey, error) {
	return "", fmt.Errorf("not object key, but unit")
}
func (u Unit) IsArrayIndex() bool { return false }
func (u Unit) AsArrayIndex() (ArrayIndex, error) {
	return 0, fmt.Errorf("not array index, but unit")
}
