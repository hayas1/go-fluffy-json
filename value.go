package fluffyjson

import (
	"encoding/json"
	"fmt"
)

const (
	ObjectType = "object"
	ArrayType  = "array"
	StringType = "string"
)

type (
	JsonValue interface {
		json.Unmarshaler
		Represent() string
	}

	Value  struct{ Value JsonValue }
	Object map[string]JsonValue // TODO int key
	Array  []JsonValue
	String string
)

func New(v interface{}) (JsonValue, error) {
	var err error
	switch t := v.(type) {
	case map[string]interface{}:
		inner := make(map[string]JsonValue, len(t))
		for k, v := range t {
			if inner[k], err = New(v); err != nil {
				return nil, err
			}
		}
		return &[]Object{(inner)}[0], nil
	case []interface{}:
		inner := make([]JsonValue, len(t))
		for i, v := range t {
			if inner[i], err = New(v); err != nil {
				return nil, err
			}
		}
		return &[]Array{(inner)}[0], nil
	case string:
		return &[]String{String(t)}[0], nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}

func (v Value) Represent() string {
	return v.Value.Represent()
}
func (v *Value) UnmarshalJSON(data []byte) error {
	// TODO remove this wrapper struct `Value` ?
	// TODO do not implement as deep copy, unmarshal directly
	var inner interface{}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	} else if value, err := New(inner); err != nil {
		return err
	} else {
		v.Value = value
	}
	return nil
}

func (o Object) Represent() string {
	return ObjectType
}
func (o *Object) UnmarshalJSON(data []byte) error {
	var inner interface{}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	} else if object, err := New(inner); err != nil {
		return err
	} else {
		*o = *object.(*Object)
	}
	return nil
}

func (a Array) Represent() string {
	return ArrayType
}
func (a *Array) UnmarshalJSON(data []byte) error {
	var inner interface{}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	} else if array, err := New(inner); err != nil {
		return err
	} else {
		*a = *array.(*Array)
	}
	return nil
}

func (s String) Represent() string {
	return StringType
}

func (s *String) UnmarshalJSON(data []byte) error {
	var inner interface{}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	str, err := New(inner)
	if err != nil {
		return err
	}
	*s = *str.(*String)
	return nil
}
