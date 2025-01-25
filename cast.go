package fluffyjson

import (
	"fmt"
)

type (
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
	ErrAsValue struct {
		Not Representation
		But Representation
	}
	ErrCast struct {
		Unsupported interface{}
	}
)

func (e ErrAsValue) Error() string {
	return fmt.Sprintf("not %s, but %s", e.Not, e.But)
}
func (e ErrCast) Error() string {
	return fmt.Sprintf("unsupported type %T", e.Unsupported)
}

func Cast(a interface{}) (JsonValue, error) {
	switch t := a.(type) {
	case map[string]interface{}:
		o, err := CastObject(t)
		return &o, err
	case []interface{}:
		a, err := CastArray(t)
		return &a, err
	case string:
		s, err := CastString(t)
		return &s, err
	default:
		return nil, ErrCast{Unsupported: t}
	}
}
func (v Value) IsObject() bool            { return v.Value.IsObject() }
func (v Value) AsObject() (Object, error) { return v.Value.AsObject() }
func (v Value) IsArray() bool             { return v.Value.IsArray() }
func (v Value) AsArray() (Array, error)   { return v.Value.AsArray() }
func (v Value) IsString() bool            { return v.Value.IsString() }
func (v Value) AsString() (String, error) { return v.Value.AsString() }

func CastObject(m map[string]interface{}) (Object, error) {
	var err error
	object := make(map[string]JsonValue, len(m))
	for k, v := range m {
		if object[k], err = Cast(v); err != nil {
			return nil, err
		}
	}
	return object, nil
}
func (o Object) IsObject() bool            { return true }
func (o Object) AsObject() (Object, error) { return o, nil }
func (o Object) IsArray() bool             { return false }
func (o Object) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: OBJECT} }
func (o Object) IsString() bool            { return false }
func (o Object) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: OBJECT} }

func CastArray(l []interface{}) (Array, error) {
	var err error
	array := make([]JsonValue, len(l))
	for i, v := range l {
		if array[i], err = Cast(v); err != nil {
			return nil, err
		}
	}
	return array, nil
}
func (a Array) IsObject() bool            { return false }
func (a Array) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: ARRAY} }
func (a Array) IsArray() bool             { return true }
func (a Array) AsArray() (Array, error)   { return a, nil }
func (a Array) IsString() bool            { return false }
func (a Array) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: ARRAY} }

func CastString(s string) (String, error) {
	return String(s), nil
}
func (s String) IsObject() bool            { return false }
func (s String) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: STRING} }
func (s String) IsArray() bool             { return false }
func (s String) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: STRING} }
func (s String) IsString() bool            { return true }
func (s String) AsString() (String, error) { return s, nil }
