package fluffyjson

import (
	"iter"
)

type (
	Accept interface {
		Accept(Visitor) error
	}
	Visitor interface {
		VisitRoot(*Value) error
		LeaveRoot(*Value) error

		VisitObject(*Object) error
		VisitObjectEntry(string, JsonValue) error
		LeaveObjectEntry(string, JsonValue) error
		LeaveObject(*Object) error

		VisitArray(*Array) error
		VisitArrayEntry(int, JsonValue) error
		LeaveArrayEntry(int, JsonValue) error
		LeaveArray(*Array) error

		VisitString(*String) error
		VisitNumber(*Number) error
		VisitBool(*Bool) error
		VisitNull(*Null) error
	}

	BaseVisitor    struct{}
	Dfs[V Visitor] struct {
		Visitor V
	}
)

func (v *Value) Accept(visitor Visitor) error  { return visitor.VisitRoot(v) }
func (o *Object) Accept(visitor Visitor) error { return visitor.VisitObject(o) }
func (a *Array) Accept(visitor Visitor) error  { return visitor.VisitArray(a) }
func (s *String) Accept(visitor Visitor) error { return visitor.VisitString(s) }
func (n *Number) Accept(visitor Visitor) error { return visitor.VisitNumber(n) }
func (b *Bool) Accept(visitor Visitor) error   { return visitor.VisitBool(b) }
func (n *Null) Accept(visitor Visitor) error   { return visitor.VisitNull(n) }

func (bv *BaseVisitor) VisitRoot(v *Value) error                         { return nil }
func (bv *BaseVisitor) LeaveRoot(v *Value) error                         { return nil }
func (bv *BaseVisitor) VisitObject(o *Object) error                      { return nil }
func (bv *BaseVisitor) VisitObjectEntry(key string, val JsonValue) error { return nil }
func (bv *BaseVisitor) LeaveObjectEntry(key string, val JsonValue) error { return nil }
func (bv *BaseVisitor) LeaveObject(o *Object) error                      { return nil }
func (bv *BaseVisitor) VisitArray(a *Array) error                        { return nil }
func (bv *BaseVisitor) VisitArrayEntry(idx int, val JsonValue) error     { return nil }
func (bv *BaseVisitor) LeaveArrayEntry(idx int, val JsonValue) error     { return nil }
func (bv *BaseVisitor) LeaveArray(a *Array) error                        { return nil }
func (bv *BaseVisitor) VisitString(s *String) error                      { return nil }
func (bv *BaseVisitor) VisitNumber(n *Number) error                      { return nil }
func (bv *BaseVisitor) VisitBool(b *Bool) error                          { return nil }
func (bv *BaseVisitor) VisitNull(n *Null) error                          { return nil }

func DfsVisitor[V Visitor](visitor V) *Dfs[V] {
	return &Dfs[V]{Visitor: visitor}
}
func (dfs *Dfs[V]) VisitRoot(v *Value) (err error) {
	if err = dfs.Visitor.VisitRoot(v); err != nil {
		return err
	}
	if err := v.Value.Accept(dfs); err != nil {
		return err
	}
	defer func() { err = dfs.LeaveRoot(v) }() // TODO if err != nil, leave should be called or not ?
	return nil
}
func (dfs *Dfs[V]) LeaveRoot(v *Value) error {
	return dfs.Visitor.LeaveRoot(v)
}
func (dfs *Dfs[V]) VisitObject(o *Object) (err error) {
	if err = dfs.Visitor.VisitObject(o); err != nil {
		return err
	}
	defer func() { err = dfs.LeaveObject(o) }()

	for k, v := range *o {
		if err := dfs.VisitObjectEntry(k, v); err != nil {
			return err
		}
	}
	return nil
}
func (dfs *Dfs[V]) VisitObjectEntry(k string, v JsonValue) (err error) {
	if err = dfs.Visitor.VisitObjectEntry(k, v); err != nil {
		return err
	}
	defer func() { err = dfs.LeaveObjectEntry(k, v) }()

	if err := v.Accept(dfs); err != nil {
		return err
	}
	return nil
}
func (dfs *Dfs[V]) LeaveObjectEntry(k string, v JsonValue) error {
	return dfs.Visitor.LeaveObjectEntry(k, v)
}
func (dfs *Dfs[V]) LeaveObject(o *Object) error {
	return dfs.Visitor.LeaveObject(o)
}
func (dfs *Dfs[V]) VisitArray(a *Array) (err error) {
	if err = dfs.Visitor.VisitArray(a); err != nil {
		return err
	}
	defer func() { err = dfs.LeaveArray(a) }()

	for i, v := range *a {
		if err := dfs.VisitArrayEntry(i, v); err != nil {
			return err
		}
	}
	return nil
}
func (dfs *Dfs[V]) VisitArrayEntry(i int, v JsonValue) (err error) {
	if err = dfs.Visitor.VisitArrayEntry(i, v); err != nil {
		return err
	}
	defer func() { err = dfs.LeaveArrayEntry(i, v) }()

	if err := v.Accept(dfs); err != nil {
		return err
	}
	return nil
}
func (dfs *Dfs[V]) LeaveArrayEntry(i int, v JsonValue) error {
	return dfs.Visitor.LeaveArrayEntry(i, v)
}
func (dfs *Dfs[V]) LeaveArray(a *Array) error {
	return dfs.Visitor.LeaveArray(a)
}
func (dfs *Dfs[V]) VisitString(s *String) error {
	return dfs.Visitor.VisitString(s)
}
func (dfs *Dfs[V]) VisitNumber(n *Number) error {
	return dfs.Visitor.VisitNumber(n)
}
func (dfs *Dfs[V]) VisitBool(b *Bool) error {
	return dfs.Visitor.VisitBool(b)
}
func (dfs *Dfs[V]) VisitNull(n *Null) error {
	return dfs.Visitor.VisitNull(n)
}

type ValueVisitor struct {
	BaseVisitor
	pointer Pointer
	yield   func(Pointer, JsonValue) bool
}

func (v *Value) DepthFirst() iter.Seq2[Pointer, JsonValue] {
	return func(yield func(Pointer, JsonValue) bool) {
		visitor := &ValueVisitor{yield: yield}
		v.Accept(DfsVisitor(visitor))
	}
}

func (vv *ValueVisitor) VisitRoot(v *Value) error {
	vv.pointer = nil
	vv.yield(vv.pointer, v.Value)
	return nil
}
func (vv *ValueVisitor) VisitObjectEntry(k string, v JsonValue) error {
	vv.pointer = append(vv.pointer, KeyAccess(k))
	vv.yield(vv.pointer, v)
	return nil
}
func (vv *ValueVisitor) LeaveObjectEntry(k string, v JsonValue) error {
	vv.pointer = vv.pointer[:len(vv.pointer)-1]
	return nil
}
func (vv *ValueVisitor) VisitArrayEntry(i int, v JsonValue) error {
	vv.pointer = append(vv.pointer, IndexAccess(i))
	vv.yield(vv.pointer, v)
	return nil
}
func (vv *ValueVisitor) LeaveArrayEntry(i int, v JsonValue) error {
	vv.pointer = vv.pointer[:len(vv.pointer)-1]
	return nil
}
