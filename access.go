package fluffyjson

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	Access interface {
		Access(...Accessor) (JsonValue, error)
		Slice(SliceAccessor) ([]JsonValue, error)
	}

	Accessor interface {
		Accessing(JsonValue) (JsonValue, error)
	}
	SliceAccessor interface {
		Slicing(JsonValue) ([]JsonValue, error)
	}

	KeyAccess      string
	IndexAccess    int
	KeyIndexAccess string
	Pointer        []Accessor

	SliceAccess struct {
		Start int
		End   int
	}

	ErrAccess struct {
		Accessor string
		Expected Representation
		Actual   Representation
	}
)

func (e ErrAccess) Error() string {
	return fmt.Sprintf("%s only allowed on %s, got %s", e.Accessor, e.Expected, e.Actual)
}

func (o *Object) Access(ptr ...Accessor) (JsonValue, error)    { return Pointer(ptr).Accessing(o) }
func (o *Object) Slice(acc SliceAccessor) ([]JsonValue, error) { return acc.Slicing(o) }
func (a *Array) Access(ptr ...Accessor) (JsonValue, error)     { return Pointer(ptr).Accessing(a) }
func (a *Array) Slice(acc SliceAccessor) ([]JsonValue, error)  { return acc.Slicing(a) }
func (s *String) Access(ptr ...Accessor) (JsonValue, error)    { return Pointer(ptr).Accessing(s) }
func (s *String) Slice(acc SliceAccessor) ([]JsonValue, error) { return acc.Slicing(s) }
func (n *Number) Access(ptr ...Accessor) (JsonValue, error)    { return Pointer(ptr).Accessing(n) }
func (n *Number) Slice(acc SliceAccessor) ([]JsonValue, error) { return acc.Slicing(n) }
func (b *Bool) Access(ptr ...Accessor) (JsonValue, error)      { return Pointer(ptr).Accessing(b) }
func (b *Bool) Slice(acc SliceAccessor) ([]JsonValue, error)   { return acc.Slicing(b) }
func (n *Null) Access(ptr ...Accessor) (JsonValue, error)      { return Pointer(ptr).Accessing(n) }
func (n *Null) Slice(acc SliceAccessor) ([]JsonValue, error)   { return acc.Slicing(n) }

func (k KeyAccess) Accessing(v JsonValue) (JsonValue, error) {
	switch o := v.(type) {
	case *Object:
		return (*o)[string(k)], nil
	default:
		return nil, ErrAccess{
			Accessor: fmt.Sprintf("%T", k),
			Expected: OBJECT,
			Actual:   v.Representation(),
		}
	}
}
func (i IndexAccess) Accessing(v JsonValue) (JsonValue, error) {
	switch a := v.(type) {
	case *Array:
		return (*a)[i], nil
	default:
		return nil, ErrAccess{
			Accessor: fmt.Sprintf("%T", i),
			Expected: ARRAY,
			Actual:   v.Representation(),
		}
	}
}
func (ki KeyIndexAccess) Accessing(v JsonValue) (JsonValue, error) {
	switch t := v.(type) {
	case *Object:
		return KeyAccess(ki).Accessing(t)
	case *Array:
		index, err := strconv.Atoi(string(ki))
		if err != nil {
			return nil, err
		}
		return IndexAccess(index).Accessing(t)
	default:
		return nil, ErrAccess{
			Accessor: fmt.Sprintf("%T", ki),
			Expected: OBJECT, // TODO OBJECT or ARRAY
			Actual:   v.Representation(),
		}
	}
}

func (s SliceAccess) Slicing(v JsonValue) ([]JsonValue, error) {
	switch a := v.(type) {
	case *Array:
		return (*a)[s.Start:s.End], nil
	default:
		return nil, ErrAccess{
			Accessor: fmt.Sprintf("%T", s),
			Expected: ARRAY,
			Actual:   v.Representation(),
		}
	}
}

// https://tools.ietf.org/html/rfc6901
func ParsePointer(p string) (Pointer, error) {
	if !strings.HasPrefix(p, "/") {
		return nil, fmt.Errorf("%s is not prefixed by /", p)
	} else if p == "/" {
		return nil, nil
	}

	parsed := make([]string, 0)
	for _, s := range strings.Split(p, "/")[1:] {
		s = strings.ReplaceAll(s, "~1", "/")
		s = strings.ReplaceAll(s, "~0", "~")
		parsed = append(parsed, s)
	}

	pointer := make([]Accessor, 0, len(parsed))
	for _, a := range parsed {
		pointer = append(pointer, KeyIndexAccess(a))
	}
	return pointer, nil
}
func (p Pointer) String() (string, error) {
	escaped := make([]string, 0, len(p))
	for _, acc := range p {
		var pointer string
		switch ki := acc.(type) {
		case KeyAccess:
			pointer = string(ki)
		case IndexAccess:
			pointer = fmt.Sprint(ki)
		case KeyIndexAccess:
			pointer = string(ki)
		// TODO case Pointer:
		default:
			return "", fmt.Errorf("unknown accessor %T", acc)
		}
		pointer = strings.ReplaceAll(pointer, "~", "~0")
		pointer = strings.ReplaceAll(pointer, "/", "~1")
		escaped = append(escaped, pointer)
	}
	return "/" + strings.Join(escaped, "/"), nil
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
