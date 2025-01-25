package fluffyjson

import "fmt"

type (
	Access interface {
		Access(Accessor) (JsonValue, error)
		Slice(SliceAccessor) ([]JsonValue, error)
	}

	Accessor interface {
		Accessing(JsonValue) (JsonValue, error)
	}
	SliceAccessor interface {
		Slicing(JsonValue) ([]JsonValue, error)
	}

	KeyAccess   string
	IndexAccess int
	SliceAccess struct {
		Start int
		End   int
	}

	Pointer []Accessor
)

func (v *Value) Access(accessor Accessor) (JsonValue, error)        { return accessor.Accessing(v) }
func (v *Value) Slice(accessor SliceAccessor) ([]JsonValue, error)  { return accessor.Slicing(v) }
func (o *Object) Access(accessor Accessor) (JsonValue, error)       { return accessor.Accessing(o) }
func (o *Object) Slice(accessor SliceAccessor) ([]JsonValue, error) { return accessor.Slicing(o) }
func (a *Array) Access(accessor Accessor) (JsonValue, error)        { return accessor.Accessing(a) }
func (a *Array) Slice(accessor SliceAccessor) ([]JsonValue, error)  { return accessor.Slicing(a) }
func (s *String) Access(accessor Accessor) (JsonValue, error)       { return accessor.Accessing(s) }
func (s *String) Slice(accessor SliceAccessor) ([]JsonValue, error) { return accessor.Slicing(s) }

func (k KeyAccess) Accessing(v JsonValue) (JsonValue, error) {
	switch o := v.(type) {
	case *Value:
		return k.Accessing(o.Value)
	case *Object:
		return (*o)[string(k)], nil
	default:
		return nil, fmt.Errorf("key access only allowed on object, got %T", v)
	}
}
func (i IndexAccess) Accessing(v JsonValue) (JsonValue, error) {
	switch a := v.(type) {
	case *Value:
		return i.Accessing(a.Value)
	case *Array:
		return (*a)[i], nil
	default:
		return nil, fmt.Errorf("index access only allowed on array, got %T", v)
	}
}
func (s SliceAccess) Slicing(v JsonValue) ([]JsonValue, error) {
	switch a := v.(type) {
	case *Value:
		return s.Slicing(a.Value)
	case *Array:
		return (*a)[s.Start:s.End], nil
	default:
		return nil, fmt.Errorf("slice access only allowed on array, got %T", v)
	}
}
