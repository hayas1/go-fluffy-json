package fluffyjson

import (
	"fmt"
	"strconv"
	"strings"
)

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

	ErrAccess struct {
		Accessor string
		Expect   Representation
		Actual   Representation
	}
)

func (e ErrAccess) Error() string {
	return fmt.Sprintf("%s only allowed on %s, got %s", e.Accessor, e.Expect, e.Actual)
}

func (v *RootValue) Access(acc Accessor) (JsonValue, error)       { return acc.Accessing(v) }
func (v *RootValue) Slice(acc SliceAccessor) ([]JsonValue, error) { return acc.Slicing(v) }
func (v *RootValue) Pointer(ptr ...Accessor) (JsonValue, error)   { return Pointer(ptr).Accessing(v) }
func (o *Object) Access(acc Accessor) (JsonValue, error)          { return acc.Accessing(o) }
func (o *Object) Slice(acc SliceAccessor) ([]JsonValue, error)    { return acc.Slicing(o) }
func (o *Object) Pointer(ptr ...Accessor) (JsonValue, error)      { return Pointer(ptr).Accessing(o) }
func (a *Array) Access(acc Accessor) (JsonValue, error)           { return acc.Accessing(a) }
func (a *Array) Slice(acc SliceAccessor) ([]JsonValue, error)     { return acc.Slicing(a) }
func (a *Array) Pointer(ptr ...Accessor) (JsonValue, error)       { return Pointer(ptr).Accessing(a) }
func (s *String) Access(acc Accessor) (JsonValue, error)          { return acc.Accessing(s) }
func (s *String) Slice(acc SliceAccessor) ([]JsonValue, error)    { return acc.Slicing(s) }
func (s *String) Pointer(ptr ...Accessor) (JsonValue, error)      { return Pointer(ptr).Accessing(s) }
func (n *Number) Access(acc Accessor) (JsonValue, error)          { return acc.Accessing(n) }
func (n *Number) Slice(acc SliceAccessor) ([]JsonValue, error)    { return acc.Slicing(n) }
func (n *Number) Pointer(ptr ...Accessor) (JsonValue, error)      { return Pointer(ptr).Accessing(n) }
func (b *Bool) Access(acc Accessor) (JsonValue, error)            { return acc.Accessing(b) }
func (b *Bool) Slice(acc SliceAccessor) ([]JsonValue, error)      { return acc.Slicing(b) }
func (b *Bool) Pointer(ptr ...Accessor) (JsonValue, error)        { return Pointer(ptr).Accessing(b) }
func (n *Null) Access(acc Accessor) (JsonValue, error)            { return acc.Accessing(n) }
func (n *Null) Slice(acc SliceAccessor) ([]JsonValue, error)      { return acc.Slicing(n) }
func (n *Null) Pointer(ptr ...Accessor) (JsonValue, error)        { return Pointer(ptr).Accessing(n) }

func (k KeyAccess) Accessing(v JsonValue) (JsonValue, error) {
	switch o := v.(type) {
	case *RootValue:
		return k.Accessing(o.Value)
	case *Object:
		return (*o)[string(k)], nil
	default:
		return nil, ErrAccess{
			Accessor: fmt.Sprintf("%T", k),
			Expect:   OBJECT,
			Actual:   v.Representation(),
		}
	}
}
func (i IndexAccess) Accessing(v JsonValue) (JsonValue, error) {
	switch a := v.(type) {
	case *RootValue:
		return i.Accessing(a.Value)
	case *Array:
		return (*a)[i], nil
	default:
		return nil, ErrAccess{
			Accessor: fmt.Sprintf("%T", i),
			Expect:   ARRAY,
			Actual:   v.Representation(),
		}
	}
}
func (s SliceAccess) Slicing(v JsonValue) ([]JsonValue, error) {
	switch a := v.(type) {
	case *RootValue:
		return s.Slicing(a.Value)
	case *Array:
		return (*a)[s.Start:s.End], nil
	default:
		return nil, ErrAccess{
			Accessor: fmt.Sprintf("%T", s),
			Expect:   ARRAY,
			Actual:   v.Representation(),
		}
	}
}

// https://tools.ietf.org/html/rfc6901
func ParsePointer(p string) Pointer {
	if p == "/" {
		return nil
	}

	parsed := make([]string, 0)
	for _, s := range strings.Split(p, "/")[1:] {
		s = strings.ReplaceAll(s, "~1", "/")
		s = strings.ReplaceAll(s, "~0", "~")
		parsed = append(parsed, s)
	}

	pointer := make([]Accessor, 0, len(parsed))
	for _, a := range parsed {
		// TODO integer like map key
		if index, err := strconv.Atoi(a); err != nil {
			pointer = append(pointer, KeyAccess(a))
		} else {
			pointer = append(pointer, IndexAccess(index))
		}
	}
	return pointer
}
func (p Pointer) String() string {
	escaped := make([]string, 0, len(p))
	for _, acc := range p {
		var pointer string
		switch ki := acc.(type) {
		case KeyAccess:
			pointer = string(ki)
		case IndexAccess:
			pointer = fmt.Sprint(ki)
		}
		pointer = strings.ReplaceAll(pointer, "~", "~0")
		pointer = strings.ReplaceAll(pointer, "/", "~1")
		escaped = append(escaped, pointer)
	}
	return "/" + strings.Join(escaped, "/")
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
