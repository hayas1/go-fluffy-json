// [fluffyjson] can deal with JSON fluffily.
//   - Compatible with [encode/json] and better interface than interface{}.
//   - Useful methods to handle JSON value such as cast, access, visit, and so on.
//   - Pure Go implementation.
package fluffyjson

import (
	"encoding/json"
	"fmt"
)

type (
	// The interface of JSON(https://www.json.org/) value
	JsonValue interface {
		representation() representation
		json.Unmarshaler // TODO this cause pointer receiver
		json.Marshaler
		As
		Access
		AccessAs
		Accept
		Search
	}

	representation string

	RootValue struct{ JsonValue }
	Object    map[string]JsonValue
	Array     []JsonValue
	String    string
	Number    float64
	Bool      bool
	Null      func(null)

	null struct {
		_ struct{}
	}

	ErrCast struct {
		Unsupported any
	}
	ErrUnmarshal struct {
		Data []byte
	}
)

func (e ErrCast) Error() string {
	return fmt.Sprintf("unsupported type %T", e.Unsupported)
}
func (e ErrUnmarshal) Error() string {
	return fmt.Sprintf("cannot unmarshal %s", e.Data)
}

const (
	OBJECT representation = "object"
	ARRAY  representation = "array"
	STRING representation = "string"
	NUMBER representation = "number"
	BOOL   representation = "bool"
	NULL   representation = "null"
)

func (v *RootValue) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case '{':
		var o Object
		if err := json.Unmarshal(data, &o); err != nil {
			return err
		}
		v.JsonValue = &o
	case '[':
		var a Array
		if err := json.Unmarshal(data, &a); err != nil {
			return err
		}
		v.JsonValue = &a
	case '"':
		var s String
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		v.JsonValue = &s
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var n Number
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		v.JsonValue = &n
	case 't', 'f':
		var b Bool
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}
		v.JsonValue = &b
	case 'n':
		var n Null
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		v.JsonValue = &n
	default:
		var unmarshaled interface{}
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			return err
		}
		value, err := Cast(unmarshaled)
		if err != nil {
			return err
		}
		v.JsonValue = value
		return ErrUnmarshal{Data: data}
	}
	return nil
}
func (v RootValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.JsonValue)
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

func (o Object) representation() representation { return OBJECT }
func (o *Object) UnmarshalJSON(data []byte) error {
	var inner any
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

func (a Array) representation() representation { return ARRAY }
func (a *Array) UnmarshalJSON(data []byte) error {
	var inner any
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

func (s *String) representation() representation { return STRING }
func (s *String) UnmarshalJSON(data []byte) error {
	var inner any
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
func CastString(s string) (String, error) {
	return String(s), nil
}

func (n *Number) representation() representation { return NUMBER }
func (n *Number) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	num, err := Cast(inner)
	if err != nil {
		return err
	}
	*n = *num.(*Number)
	return nil
}
func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(n))
}
func CastNumber(n float64) (Number, error) {
	return Number(n), nil
}

func (b *Bool) representation() representation { return BOOL }
func (b *Bool) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	bool, err := Cast(inner)
	if err != nil {
		return err
	}
	*b = *bool.(*Bool)
	return nil
}
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(b))
}
func CastBool(b bool) (Bool, error) {
	return Bool(b), nil
}

func (n *Null) representation() representation { return NULL }
func (n *Null) UnmarshalJSON(data []byte) error {
	var inner any
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	null, err := Cast(inner)
	if err != nil {
		return err
	}
	*n = *null.(*Null)
	return nil
}
func (n Null) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
func CastNull(n func(null)) (Null, error) {
	return Null(n), nil
}
