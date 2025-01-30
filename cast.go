package fluffyjson

import "fmt"

type (
	As interface {
		AsObject
		AsArray
		AsString
		AsNumber
		AsBool
		AsNull
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
		Expected Representation
		Actual   Representation
	}
)

func (e ErrAsValue) Error() string {
	return fmt.Sprintf("not %s, but %s", e.Expected, e.Actual)
}

func Cast(v any) (JsonValue, error) {
	switch t := v.(type) {
	case map[string]any:
		o, err := CastObject(t)
		return &o, err
	case []any:
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

func (o Object) IsObject() bool            { return true }
func (o Object) AsObject() (Object, error) { return o, nil }
func (o Object) IsArray() bool             { return false }
func (o Object) AsArray() (Array, error)   { return nil, ErrAsValue{Expected: ARRAY, Actual: OBJECT} }
func (o Object) IsString() bool            { return false }
func (o Object) AsString() (String, error) { return "", ErrAsValue{Expected: STRING, Actual: OBJECT} }
func (o Object) IsNumber() bool            { return false }
func (o Object) AsNumber() (Number, error) { return 0, ErrAsValue{Expected: NUMBER, Actual: OBJECT} }
func (o Object) IsBool() bool              { return false }
func (o Object) AsBool() (Bool, error)     { return false, ErrAsValue{Expected: BOOL, Actual: OBJECT} }
func (o Object) IsNull() bool              { return false }
func (o Object) AsNull() (Null, error)     { return nil, ErrAsValue{Expected: NULL, Actual: OBJECT} }

func (a Array) IsObject() bool            { return false }
func (a Array) AsObject() (Object, error) { return nil, ErrAsValue{Expected: OBJECT, Actual: ARRAY} }
func (a Array) IsArray() bool             { return true }
func (a Array) AsArray() (Array, error)   { return a, nil }
func (a Array) IsString() bool            { return false }
func (a Array) AsString() (String, error) { return "", ErrAsValue{Expected: STRING, Actual: ARRAY} }
func (a Array) IsNumber() bool            { return false }
func (a Array) AsNumber() (Number, error) { return 0, ErrAsValue{Expected: NUMBER, Actual: ARRAY} }
func (a Array) IsBool() bool              { return false }
func (a Array) AsBool() (Bool, error)     { return false, ErrAsValue{Expected: BOOL, Actual: ARRAY} }
func (a Array) IsNull() bool              { return false }
func (a Array) AsNull() (Null, error)     { return nil, ErrAsValue{Expected: NULL, Actual: ARRAY} }

func (s String) IsArray() bool             { return false }
func (s String) AsArray() (Array, error)   { return nil, ErrAsValue{Expected: ARRAY, Actual: STRING} }
func (s String) IsObject() bool            { return false }
func (s String) AsObject() (Object, error) { return nil, ErrAsValue{Expected: OBJECT, Actual: STRING} }
func (s String) IsString() bool            { return true }
func (s String) AsString() (String, error) { return s, nil }
func (s String) IsNumber() bool            { return false }
func (s String) AsNumber() (Number, error) { return 0, ErrAsValue{Expected: NUMBER, Actual: STRING} }
func (s String) IsBool() bool              { return false }
func (s String) AsBool() (Bool, error)     { return false, ErrAsValue{Expected: BOOL, Actual: STRING} }
func (s String) IsNull() bool              { return false }
func (s String) AsNull() (Null, error)     { return nil, ErrAsValue{Expected: NULL, Actual: STRING} }

func (n Number) IsObject() bool            { return false }
func (n Number) AsObject() (Object, error) { return nil, ErrAsValue{Expected: OBJECT, Actual: NUMBER} }
func (n Number) IsArray() bool             { return false }
func (n Number) AsArray() (Array, error)   { return nil, ErrAsValue{Expected: ARRAY, Actual: NUMBER} }
func (n Number) IsString() bool            { return false }
func (n Number) AsString() (String, error) { return "", ErrAsValue{Expected: STRING, Actual: NUMBER} }
func (n Number) IsNumber() bool            { return true }
func (n Number) AsNumber() (Number, error) { return n, nil }
func (n Number) IsBool() bool              { return false }
func (n Number) AsBool() (Bool, error)     { return false, ErrAsValue{Expected: BOOL, Actual: NUMBER} }
func (n Number) IsNull() bool              { return false }
func (n Number) AsNull() (Null, error)     { return nil, ErrAsValue{Expected: NULL, Actual: NUMBER} }

func (b Bool) IsObject() bool            { return false }
func (b Bool) AsObject() (Object, error) { return nil, ErrAsValue{Expected: OBJECT, Actual: BOOL} }
func (b Bool) IsArray() bool             { return false }
func (b Bool) AsArray() (Array, error)   { return nil, ErrAsValue{Expected: ARRAY, Actual: BOOL} }
func (b Bool) IsString() bool            { return false }
func (b Bool) AsString() (String, error) { return "", ErrAsValue{Expected: STRING, Actual: BOOL} }
func (b Bool) IsNumber() bool            { return false }
func (b Bool) AsNumber() (Number, error) { return 0, ErrAsValue{Expected: NUMBER, Actual: BOOL} }
func (b Bool) IsBool() bool              { return true }
func (b Bool) AsBool() (Bool, error)     { return b, nil }
func (b Bool) IsNull() bool              { return false }
func (b Bool) AsNull() (Null, error)     { return nil, ErrAsValue{Expected: NULL, Actual: BOOL} }

func (n Null) IsObject() bool            { return false }
func (n Null) AsObject() (Object, error) { return nil, ErrAsValue{Expected: OBJECT, Actual: NULL} }
func (n Null) IsArray() bool             { return false }
func (n Null) AsArray() (Array, error)   { return nil, ErrAsValue{Expected: ARRAY, Actual: NULL} }
func (n Null) IsString() bool            { return false }
func (n Null) AsString() (String, error) { return "", ErrAsValue{Expected: STRING, Actual: NULL} }
func (n Null) IsNumber() bool            { return false }
func (n Null) AsNumber() (Number, error) { return 0, ErrAsValue{Expected: NUMBER, Actual: NULL} }
func (n Null) IsBool() bool              { return false }
func (n Null) AsBool() (Bool, error)     { return false, ErrAsValue{Expected: BOOL, Actual: NULL} }
func (n Null) IsNull() bool              { return true }
func (n Null) AsNull() (Null, error)     { return n, nil }

func CastObject(m map[string]any) (Object, error) {
	var err error
	object := make(map[string]JsonValue, len(m))
	for k, v := range m {
		if object[k], err = Cast(v); err != nil {
			return nil, err
		}
	}
	return object, nil
}
func CastArray(l []any) (Array, error) {
	var err error
	array := make([]JsonValue, len(l))
	for i, v := range l {
		if array[i], err = Cast(v); err != nil {
			return nil, err
		}
	}
	return array, nil
}
func CastString(s string) (String, error) {
	return String(s), nil
}
func CastNumber(n float64) (Number, error) {
	return Number(n), nil
}
func CastBool(b bool) (Bool, error) {
	return Bool(b), nil
}
func CastNull(n func(null)) (Null, error) {
	return Null(n), nil
}
