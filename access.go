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

func (k KeyAccess) Access(v JsonValue) (JsonValue, error) {
	switch o := v.(type) {
	case *Object:
		return (*o)[string(k)], nil
	default:
		return nil, fmt.Errorf("key access only allowed on object, got %T", v)
	}
}
func (i IndexAccess) Access(v JsonValue) (JsonValue, error) {
	switch a := v.(type) {
	case *Array:
		return (*a)[i], nil
	default:
		return nil, fmt.Errorf("index access only allowed on array, got %T", v)
	}
}
func (s SliceAccess) SliceAccess(v JsonValue) ([]JsonValue, error) {
	switch a := v.(type) {
	case *Array:
		return (*a)[s.Start:s.End], nil
	default:
		return nil, fmt.Errorf("slice access only allowed on array, got %T", v)
	}
}
