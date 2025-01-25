package fluffyjson

import "fmt"

type (
	Access interface {
		Access(Accessor) (JsonValue, error)
		Slice(SliceAccessor) ([]JsonValue, error)
	}

	Accessor interface {
		Access(JsonValue) (JsonValue, error)
	}
	SliceAccessor interface {
		SliceAccess(JsonValue) ([]JsonValue, error)
	}

	KeyAccess   string
	IndexAccess int
	SliceAccess struct {
		Start int
		End   int
	}

	Pointer []Accessor
)

func (v *Value) Access(accessor Accessor) (JsonValue, error)        { return accessor.Access(v) }
func (v *Value) Slice(accessor SliceAccessor) ([]JsonValue, error)  { return accessor.SliceAccess(v) }
func (o *Object) Access(accessor Accessor) (JsonValue, error)       { return accessor.Access(o) }
func (o *Object) Slice(accessor SliceAccessor) ([]JsonValue, error) { return accessor.SliceAccess(o) }
func (a *Array) Access(accessor Accessor) (JsonValue, error)        { return accessor.Access(a) }
func (a *Array) Slice(accessor SliceAccessor) ([]JsonValue, error)  { return accessor.SliceAccess(a) }
func (s *String) Access(accessor Accessor) (JsonValue, error)       { return accessor.Access(s) }
func (s *String) Slice(accessor SliceAccessor) ([]JsonValue, error) { return accessor.SliceAccess(s) }

func (k KeyAccess) Access(v JsonValue) (JsonValue, error) {
	switch o := v.(type) {
	case *Value:
		return k.Access(o.Value)
	case *Object:
		return (*o)[string(k)], nil
	default:
		return nil, fmt.Errorf("key access only allowed on object, got %T", v)
	}
}
func (i IndexAccess) Access(v JsonValue) (JsonValue, error) {
	switch a := v.(type) {
	case *Value:
		return i.Access(a.Value)
	case *Array:
		return (*a)[i], nil
	default:
		return nil, fmt.Errorf("index access only allowed on array, got %T", v)
	}
}
func (s SliceAccess) SliceAccess(v JsonValue) ([]JsonValue, error) {
	switch a := v.(type) {
	case *Value:
		return s.SliceAccess(a.Value)
	case *Array:
		return (*a)[s.Start:s.End], nil
	default:
		return nil, fmt.Errorf("slice access only allowed on array, got %T", v)
	}
}
