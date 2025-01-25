package fluffyjson

import (
	"encoding/json"
	"fmt"
)

type (
	JsonValue interface {
		json.Unmarshaler
		json.Marshaler
		AsObject
		AsArray
		AsString
		Shallow
	}
	AsObject interface {
		IsObject() bool
		AsObject() (Object, error)
	}
	AsArray interface {
		IsArray() bool
		AsArray() (Array, error)
	}
	AsString interface {
		IsString() bool
		AsString() (String, error)
	}

	Value  struct{ Value JsonValue }
	Object map[string]JsonValue // TODO int key
	Array  []JsonValue
	String string
)

func Cast(v interface{}) (JsonValue, error) {
	var err error
	switch t := v.(type) {
	case map[string]interface{}:
		inner := make(map[string]JsonValue, len(t))
		for k, v := range t {
			if inner[k], err = Cast(v); err != nil {
				return nil, err
			}
		}
		return &[]Object{(inner)}[0], nil
	case []interface{}:
		inner := make([]JsonValue, len(t))
		for i, v := range t {
			if inner[i], err = Cast(v); err != nil {
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
func (v Value) IsObject() bool {
	return v.Value.IsObject()
}
func (v Value) AsObject() (Object, error) {
	return v.Value.AsObject()
}
func (v Value) IsArray() bool {
	return v.Value.IsArray()
}
func (v Value) AsArray() (Array, error) {
	return v.Value.AsArray()
}
func (v Value) IsString() bool {
	return v.Value.IsString()
}
func (v Value) AsString() (String, error) {
	return v.Value.AsString()
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
func (o Object) IsObject() bool {
	return true
}
func (o Object) AsObject() (Object, error) {
	return o, nil
}
func (o Object) IsArray() bool {
	return false
}
func (o Object) AsArray() (Array, error) {
	return nil, fmt.Errorf("not array, but object")
}
func (o Object) IsString() bool {
	return false
}
func (o Object) AsString() (String, error) {
	return "", fmt.Errorf("not string, but object")
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
func (a Array) IsObject() bool {
	return false
}
func (a Array) AsObject() (Object, error) {
	return nil, fmt.Errorf("not object, but array")
}
func (a Array) IsArray() bool {
	return true
}
func (a Array) AsArray() (Array, error) {
	return a, nil
}
func (a Array) IsString() bool {
	return false
}
func (a Array) AsString() (String, error) {
	return "", fmt.Errorf("not string, but array")
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
func (s String) IsObject() bool {
	return false
}
func (s String) AsObject() (Object, error) {
	return nil, fmt.Errorf("not object, but string")
}
func (s String) IsArray() bool {
	return false
}
func (s String) AsArray() (Array, error) {
	return nil, fmt.Errorf("not array, but string")
}
func (s String) IsString() bool {
	return true
}
func (s String) AsString() (String, error) {
	return s, nil
}
