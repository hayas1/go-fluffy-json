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
		return nil, ErrCast{Unsupported: t}
	}
}

func (e ErrAsValue) Error() string {
	return fmt.Sprintf("not %s, but %s", e.Not, e.But)
}
func (e ErrCast) Error() string {
	return fmt.Sprintf("unsupported type %T", e.Unsupported)
}

func (v Value) IsObject() bool            { return v.Value.IsObject() }
func (v Value) AsObject() (Object, error) { return v.Value.AsObject() }
func (v Value) IsArray() bool             { return v.Value.IsArray() }
func (v Value) AsArray() (Array, error)   { return v.Value.AsArray() }
func (v Value) IsString() bool            { return v.Value.IsString() }
func (v Value) AsString() (String, error) { return v.Value.AsString() }

func (o Object) IsObject() bool            { return true }
func (o Object) AsObject() (Object, error) { return o, nil }
func (o Object) IsArray() bool             { return false }
func (o Object) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: OBJECT} }
func (o Object) IsString() bool            { return false }
func (o Object) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: OBJECT} }

func (a Array) IsObject() bool            { return false }
func (a Array) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: ARRAY} }
func (a Array) IsArray() bool             { return true }
func (a Array) AsArray() (Array, error)   { return a, nil }
func (a Array) IsString() bool            { return false }
func (a Array) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: ARRAY} }

func (s String) IsObject() bool            { return false }
func (s String) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: STRING} }
func (s String) IsArray() bool             { return false }
func (s String) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: STRING} }
func (s String) IsString() bool            { return true }
func (s String) AsString() (String, error) { return s, nil }
