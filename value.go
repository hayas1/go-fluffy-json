// [fluffyjson] can deal with JSON fluffily.
//   - Compatible with [encode/json] and better interface than interface{}.
//   - Useful methods to handle JSON value such as cast, access, visit, and so on.
//   - Pure Go implementation.
package fluffyjson

import (
	"encoding/json"
)

type (
	// The interface of JSON(https://www.json.org/) value
	JsonValue interface {
		Is(Representation) bool
		Representation() Representation
		json.Unmarshaler // TODO this cause pointer receiver
		json.Marshaler
		AsValue
		Access
		AccessAs
		Accept
		Search
	}

	Representation string

	RootValue struct{ JsonValue }
	Object    map[string]JsonValue
	Array     []JsonValue
	String    string
	Number    float64
	Bool      bool
	Null      []struct{}
)

const (
	OBJECT Representation = "object"
	ARRAY  Representation = "array"
	STRING Representation = "string"
	NUMBER Representation = "number"
	BOOL   Representation = "bool"
	NULL   Representation = "null"
)

func (v *RootValue) UnmarshalJSON(data []byte) error {
	// TODO remove this wrapper struct `Value` ?
	// TODO do not implement as deep copy, unmarshal directly
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	} else if value, err := Cast(inner); err != nil {
		return err
	} else {
		v.JsonValue = value
	}
	return nil
}
func (v RootValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.JsonValue)
}

func (v Object) Is(r Representation) bool       { return r == OBJECT }
func (o Object) Representation() Representation { return OBJECT }
func (o *Object) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	} else if object, err := Cast(inner); err != nil {
		return err
	} else {
		*o = *object.(*Object)
	}
	return nil
}
func (o Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]JsonValue(o))
}

func (a Array) Is(r Representation) bool       { return r == ARRAY }
func (a Array) Representation() Representation { return ARRAY }
func (a *Array) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	} else if array, err := Cast(inner); err != nil {
		return err
	} else {
		*a = *array.(*Array)
	}
	return nil
}
func (a Array) MarshalJSON() ([]byte, error) {
	return json.Marshal([]JsonValue(a))
}

func (s String) Is(r Representation) bool       { return r == STRING }
func (s String) Representation() Representation { return STRING }
func (s *String) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	str, err := Cast(inner)
	if err != nil {
		return err
	}
	*s = *str.(*String)
	return nil
}
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (n Number) Is(r Representation) bool       { return r == NUMBER }
func (n Number) Representation() Representation { return NUMBER }
func (n *Number) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	num, err := Cast(inner)
	if err != nil {
		return err
	}
	*n = *num.(*Number)
	return nil
}
func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(n))
}

func (b Bool) Is(r Representation) bool       { return r == BOOL }
func (b Bool) Representation() Representation { return BOOL }
func (b *Bool) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	bool, err := Cast(inner)
	if err != nil {
		return err
	}
	*b = *bool.(*Bool)
	return nil
}
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(b))
}

func (n Null) Is(r Representation) bool       { return r == NULL }
func (n Null) Representation() Representation { return NULL }
func (n *Null) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	null, err := Cast(inner)
	if err != nil {
		return err
	}
	*n = *null.(*Null)
	return nil
}
func (n Null) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
