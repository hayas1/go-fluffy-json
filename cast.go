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
	AsNumber interface {
		IsNumber() bool
		AsNumber() (Number, error)
	}
	AsBool interface {
		IsBool() bool
		AsBool() (Bool, error)
	}
	AsNull interface {
		IsNull() bool
		AsNull() (Null, error)
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

func Cast(v interface{}) (JsonValue, error) {
	switch t := v.(type) {
	case map[string]interface{}:
		o, err := CastObject(t)
		return &o, err
	case []interface{}:
		a, err := CastArray(t)
		return &a, err
	case string:
		s, err := CastString(t)
		return &s, err
	case float64:
		n, err := CastNumber(t)
		return &n, err
	case bool:
		b, err := CastBool(t)
		return &b, err
	case nil:
		n, err := CastNull(nil)
		return &n, err
	default:
		return nil, ErrCast{Unsupported: t}
	}
}
func Force(v interface{}) *JsonValue {
	value, err := Cast(v)
	if err != nil {
		panic(err)
	}
	return &value
}
func (v Value) IsObject() bool            { return v.Value.IsObject() }
func (v Value) AsObject() (Object, error) { return v.Value.AsObject() }
func (v Value) IsArray() bool             { return v.Value.IsArray() }
func (v Value) AsArray() (Array, error)   { return v.Value.AsArray() }
func (v Value) IsString() bool            { return v.Value.IsString() }
func (v Value) AsString() (String, error) { return v.Value.AsString() }
func (v Value) IsNumber() bool            { return v.Value.IsNumber() }
func (v Value) AsNumber() (Number, error) { return v.Value.AsNumber() }
func (v Value) IsBool() bool              { return v.Value.IsBool() }
func (v Value) AsBool() (Bool, error)     { return v.Value.AsBool() }
func (v Value) IsNull() bool              { return v.Value.IsNull() }
func (v Value) AsNull() (Null, error)     { return v.Value.AsNull() }

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
func ForceObject(m map[string]interface{}) *Object {
	object, err := CastObject(m)
	if err != nil {
		panic(err)
	}
	return &object
}
func (o Object) IsObject() bool            { return true }
func (o Object) AsObject() (Object, error) { return o, nil }
func (o Object) IsArray() bool             { return false }
func (o Object) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: OBJECT} }
func (o Object) IsString() bool            { return false }
func (o Object) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: OBJECT} }
func (o Object) IsNumber() bool            { return false }
func (o Object) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: OBJECT} }
func (o Object) IsBool() bool              { return false }
func (o Object) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: OBJECT} }
func (o Object) IsNull() bool              { return false }
func (o Object) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: OBJECT} }

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
func ForceArray(l []interface{}) *Array {
	array, err := CastArray(l)
	if err != nil {
		panic(err)
	}
	return &array
}
func (a Array) IsObject() bool            { return false }
func (a Array) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: ARRAY} }
func (a Array) IsArray() bool             { return true }
func (a Array) AsArray() (Array, error)   { return a, nil }
func (a Array) IsString() bool            { return false }
func (a Array) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: ARRAY} }
func (a Array) IsNumber() bool            { return false }
func (a Array) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: ARRAY} }
func (a Array) IsBool() bool              { return false }
func (a Array) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: ARRAY} }
func (a Array) IsNull() bool              { return false }
func (a Array) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: ARRAY} }

func CastString(s string) (String, error) {
	return String(s), nil
}
func ForceString(s string) *String {
	str, err := CastString(s)
	if err != nil {
		panic(err)
	}
	return &str
}
func (s String) IsObject() bool            { return false }
func (s String) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: STRING} }
func (s String) IsArray() bool             { return false }
func (s String) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: STRING} }
func (s String) IsString() bool            { return true }
func (s String) AsString() (String, error) { return s, nil }
func (s String) IsNumber() bool            { return false }
func (s String) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: STRING} }
func (s String) IsBool() bool              { return false }
func (s String) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: STRING} }
func (s String) IsNull() bool              { return false }
func (s String) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: STRING} }

func CastNumber(n float64) (Number, error) {
	return Number(n), nil
}
func ForceNumber(n float64) *Number {
	num, err := CastNumber(n)
	if err != nil {
		panic(err)
	}
	return &num
}
func (n Number) IsObject() bool            { return false }
func (n Number) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: NUMBER} }
func (n Number) IsArray() bool             { return false }
func (n Number) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: NUMBER} }
func (n Number) IsString() bool            { return false }
func (n Number) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: NUMBER} }
func (n Number) IsNumber() bool            { return true }
func (n Number) AsNumber() (Number, error) { return n, nil }
func (n Number) IsBool() bool              { return false }
func (n Number) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: NUMBER} }
func (n Number) IsNull() bool              { return false }
func (n Number) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: NUMBER} }

func CastBool(b bool) (Bool, error) {
	return Bool(b), nil
}
func ForceBool(b bool) *Bool {
	bool, err := CastBool(b)
	if err != nil {
		panic(err)
	}
	return &bool
}
func (b Bool) IsObject() bool            { return false }
func (b Bool) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: BOOL} }
func (b Bool) IsArray() bool             { return false }
func (b Bool) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: BOOL} }
func (b Bool) IsString() bool            { return false }
func (b Bool) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: BOOL} }
func (b Bool) IsNumber() bool            { return false }
func (b Bool) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: BOOL} }
func (b Bool) IsBool() bool              { return true }
func (b Bool) AsBool() (Bool, error)     { return b, nil }
func (b Bool) IsNull() bool              { return false }
func (b Bool) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: BOOL} }

func CastNull(n []struct{}) (Null, error) {
	return Null(n), nil
}
func ForceNull(n []struct{}) *Null {
	null, err := CastNull(n)
	if err != nil {
		panic(err)
	}
	return &null
}
func (n Null) IsObject() bool            { return false }
func (n Null) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: NULL} }
func (n Null) IsArray() bool             { return false }
func (n Null) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: NULL} }
func (n Null) IsString() bool            { return false }
func (n Null) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: NULL} }
func (n Null) IsNumber() bool            { return false }
func (n Null) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: NULL} }
func (n Null) IsBool() bool              { return false }
func (n Null) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: NULL} }
func (n Null) IsNull() bool              { return true }
func (n Null) AsNull() (Null, error)     { return n, nil }
