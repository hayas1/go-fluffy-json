package fluffyjson

import "fmt"

type (
	Access interface {
		Access(Accessor) (JsonValue, error)
		Slice(SliceAccessor) ([]JsonValue, error)
		Pointer(...Accessor) (JsonValue, error)
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

func (v *Value) Access(acc Accessor) (JsonValue, error)        { return acc.Accessing(v) }
func (v *Value) Slice(acc SliceAccessor) ([]JsonValue, error)  { return acc.Slicing(v) }
func (v *Value) Pointer(ptr ...Accessor) (JsonValue, error)    { return Pointer(ptr).Accessing(v) }
func (o *Object) Access(acc Accessor) (JsonValue, error)       { return acc.Accessing(o) }
func (o *Object) Slice(acc SliceAccessor) ([]JsonValue, error) { return acc.Slicing(o) }
func (o *Object) Pointer(ptr ...Accessor) (JsonValue, error)   { return Pointer(ptr).Accessing(o) }
func (a *Array) Access(acc Accessor) (JsonValue, error)        { return acc.Accessing(a) }
func (a *Array) Slice(acc SliceAccessor) ([]JsonValue, error)  { return acc.Slicing(a) }
func (a *Array) Pointer(ptr ...Accessor) (JsonValue, error)    { return Pointer(ptr).Accessing(a) }
func (s *String) Access(acc Accessor) (JsonValue, error)       { return acc.Accessing(s) }
func (s *String) Slice(acc SliceAccessor) ([]JsonValue, error) { return acc.Slicing(s) }
func (s *String) Pointer(ptr ...Accessor) (JsonValue, error)   { return Pointer(ptr).Accessing(s) }

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

func (p Pointer) Accessing(v JsonValue) (JsonValue, error) {
	curr := v
	for _, a := range p {
		var err error
		curr, err = a.Accessing(curr)
		if err != nil {
			return nil, err
		}
	}
	return curr, nil
}
