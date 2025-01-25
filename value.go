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
	Object map[string]JsonValue
	Array  []JsonValue
	String string
)

func New(v interface{}) JsonValue {
	switch t := v.(type) {
	case map[string]interface{}:
		inner := make(map[string]JsonValue, len(t))
		for k, v := range t {
			inner[k] = New(v)
		}
		return &[]Object{(inner)}[0]
	case []interface{}:
		inner := make([]JsonValue, len(t))
		for i, v := range t {
			inner[i] = New(v)
		}
		return &[]Array{(inner)}[0]
	case string:
		return &[]String{String(t)}[0]
	default:
		panic("// TODO unknown type")
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
	*o = *New(inner).(*Object)
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
	*a = *New(inner).(*Array)
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
	*s = *New(inner).(*String)
	return nil
}
