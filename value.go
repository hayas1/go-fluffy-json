package fluffyjson

import (
	"encoding/json"
)

type (
	JsonValue interface {
		json.Unmarshaler
		json.Marshaler
		AsObject
		AsArray
		AsString
		Accept
		Access
	}

	Representation string

	Value  struct{ Value JsonValue }
	Object map[string]JsonValue // TODO int key
	Array  []JsonValue
	String string
)

const (
	OBJECT Representation = "object"
	ARRAY  Representation = "array"
	STRING Representation = "string"
	NUMBER Representation = "number"
	BOOL   Representation = "bool"
	NULL   Representation = "null"
)

func (v *Value) UnmarshalJSON(data []byte) error {
	// TODO remove this wrapper struct `Value` ?
	// TODO do not implement as deep copy, unmarshal directly
	var inner interface{}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	} else if value, err := Cast(inner); err != nil {
		return err
	} else {
		v.Value = value
	}
	return nil
}
func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

func (o *Object) UnmarshalJSON(data []byte) error {
	var inner interface{}
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

func (a *Array) UnmarshalJSON(data []byte) error {
	var inner interface{}
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

func (s *String) UnmarshalJSON(data []byte) error {
	var inner interface{}
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
