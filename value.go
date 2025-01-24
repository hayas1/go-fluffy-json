package jsonvalue

import (
	"encoding/json"
)

const (
	ObjectType = "object"
	ArrayType  = "array"
	StringType = "string"
)

type JsonValue interface {
	json.Unmarshaler
	Represent() string
}

type (
	Value struct {
		Value JsonValue
	}
	Object map[string]JsonValue
	Array  []JsonValue
	String string
)

func (v Value) Represent() string {
	return v.Value.Represent()
}
func (v *Value) UnmarshalJSON(data []byte) error {
	var inner interface{}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	*v = Value{Value: &[]Value{New(inner)}[0]}
	return nil
}
func New(v interface{}) Value {
	switch t := v.(type) {
	case map[string]interface{}:
		inner := make(map[string]JsonValue, len(t))
		for k, v := range t {
			inner[k] = &[]Value{New(v)}[0]
		}
		return Value{Value: &[]Object{(inner)}[0]}
	case []interface{}:
		inner := make([]JsonValue, len(t))
		for i, v := range t {
			inner[i] = &[]Value{New(v)}[0]
		}
		return Value{Value: &[]Array{(inner)}[0]}
	case string:
		return Value{Value: &[]String{String(t)}[0]}
	default:
		return Value{}
	}
}

func (o Object) Represent() string {
	return ObjectType
}
func (o *Object) UnmarshalJSON(data []byte) error {
	var inner interface{}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	*o = *New(inner).Value.(*Object)
	return nil
}

func (a Array) Represent() string {
	return ArrayType
}
func (a *Array) UnmarshalJSON(data []byte) error {
	var inner interface{}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	*a = *New(inner).Value.(*Array)
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
	*s = *New(inner).Value.(*String)
	return nil
}
