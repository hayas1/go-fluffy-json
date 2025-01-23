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

func (o Object) Represent() string {
	return ObjectType
}
func (o *Object) UnmarshalJSON(data []byte) error {
	var inner map[string]JsonValue
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	*o = inner
	return nil
}

func (a Array) Represent() string {
	return ArrayType
}

func (s String) Represent() string {
	return StringType
}

func (s *String) UnmarshalJSON(data []byte) error {
	var inner string
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	*s = String(inner)
	return nil
}
