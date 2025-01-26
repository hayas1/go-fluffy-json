package fluffyjson

import (
	"iter"
)

type (
	Accept interface {
		Accept(Visitor) error
	}
	Search interface {
		DepthFirst() iter.Seq2[Pointer, JsonValue]
	}
	Visitor interface {
		GetPointer() Pointer
		SetPointer(Pointer)
		VisitRoot(*RootValue) error
		LeaveRoot(*RootValue) error

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
	PointerVisitor struct {
		pointer Pointer
	}
	Dfs[V Visitor] struct {
		visitor V
	}
)

func (v *RootValue) Accept(visitor Visitor) error { return visitor.VisitRoot(v) }
func (o *Object) Accept(visitor Visitor) error    { return visitor.VisitObject(o) }
func (a *Array) Accept(visitor Visitor) error     { return visitor.VisitArray(a) }
func (s *String) Accept(visitor Visitor) error    { return visitor.VisitString(s) }
func (n *Number) Accept(visitor Visitor) error    { return visitor.VisitNumber(n) }
func (b *Bool) Accept(visitor Visitor) error      { return visitor.VisitBool(b) }
func (n *Null) Accept(visitor Visitor) error      { return visitor.VisitNull(n) }

func (bv *PointerVisitor) GetPointer() Pointer                              { return bv.pointer }
func (bv *PointerVisitor) SetPointer(p Pointer)                             { bv.pointer = p }
func (bv *PointerVisitor) VisitRoot(v *RootValue) error                     { return nil }
func (bv *PointerVisitor) LeaveRoot(v *RootValue) error                     { return nil }
func (bv *PointerVisitor) VisitObject(o *Object) error                      { return nil }
func (bv *PointerVisitor) VisitObjectEntry(key string, val JsonValue) error { return nil }
func (bv *PointerVisitor) LeaveObjectEntry(key string, val JsonValue) error { return nil }
func (bv *PointerVisitor) LeaveObject(o *Object) error                      { return nil }
func (bv *PointerVisitor) VisitArray(a *Array) error                        { return nil }
func (bv *PointerVisitor) VisitArrayEntry(idx int, val JsonValue) error     { return nil }
func (bv *PointerVisitor) LeaveArrayEntry(idx int, val JsonValue) error     { return nil }
func (bv *PointerVisitor) LeaveArray(a *Array) error                        { return nil }
func (bv *PointerVisitor) VisitString(s *String) error                      { return nil }
func (bv *PointerVisitor) VisitNumber(n *Number) error                      { return nil }
func (bv *PointerVisitor) VisitBool(b *Bool) error                          { return nil }
func (bv *PointerVisitor) VisitNull(n *Null) error                          { return nil }

// Get dfs wrapped visitor
func DfsVisitor[V Visitor](visitor V) *Dfs[V] {
	return &Dfs[V]{visitor: visitor}
}
func (v *Dfs[V]) GetPointer() Pointer  { return v.visitor.GetPointer() }
func (v *Dfs[V]) SetPointer(p Pointer) { v.visitor.SetPointer(p) }
func (dfs *Dfs[V]) VisitRoot(v *RootValue) (err error) {
	if err = dfs.visitor.VisitRoot(v); err != nil {
		return err
	}
	if err := v.JsonValue.Accept(dfs); err != nil {
		return err
	}
	defer func() { err = dfs.LeaveRoot(v) }() // TODO if err != nil, leave should be called or not ?
	return nil
}
func (dfs *Dfs[V]) LeaveRoot(v *RootValue) error {
	return dfs.visitor.LeaveRoot(v)
}
func (dfs *Dfs[V]) VisitObject(o *Object) (err error) {
	if err = dfs.visitor.VisitObject(o); err != nil {
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
	dfs.visitor.SetPointer(append(dfs.visitor.GetPointer(), KeyAccess(k)))
	if err = dfs.visitor.VisitObjectEntry(k, v); err != nil {
		return err
	}
	defer func() { err = dfs.LeaveObjectEntry(k, v) }()

	if err := v.Accept(dfs); err != nil {
		return err
	}
	return nil
}
func (dfs *Dfs[V]) LeaveObjectEntry(k string, v JsonValue) error {
	defer dfs.visitor.SetPointer(dfs.visitor.GetPointer()[:len(dfs.visitor.GetPointer())-1])
	return dfs.visitor.LeaveObjectEntry(k, v)
}
func (dfs *Dfs[V]) LeaveObject(o *Object) error {
	return dfs.visitor.LeaveObject(o)
}
func (dfs *Dfs[V]) VisitArray(a *Array) (err error) {
	if err = dfs.visitor.VisitArray(a); err != nil {
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
	dfs.visitor.SetPointer(append(dfs.visitor.GetPointer(), IndexAccess(i)))
	if err = dfs.visitor.VisitArrayEntry(i, v); err != nil {
		return err
	}
	defer func() { err = dfs.LeaveArrayEntry(i, v) }()

	if err := v.Accept(dfs); err != nil {
		return err
	}
	return nil
}
func (dfs *Dfs[V]) LeaveArrayEntry(i int, v JsonValue) error {
	defer dfs.visitor.SetPointer(dfs.visitor.GetPointer()[:len(dfs.visitor.GetPointer())-1])
	return dfs.visitor.LeaveArrayEntry(i, v)
}
func (dfs *Dfs[V]) LeaveArray(a *Array) error {
	return dfs.visitor.LeaveArray(a)
}
func (dfs *Dfs[V]) VisitString(s *String) error {
	return dfs.visitor.VisitString(s)
}
func (dfs *Dfs[V]) VisitNumber(n *Number) error {
	return dfs.visitor.VisitNumber(n)
}
func (dfs *Dfs[V]) VisitBool(b *Bool) error {
	return dfs.visitor.VisitBool(b)
}
func (dfs *Dfs[V]) VisitNull(n *Null) error {
	return dfs.visitor.VisitNull(n)
}

type ValueVisitor struct {
	PointerVisitor
	yield func(Pointer, JsonValue) bool
}

func depthFirstValues(v JsonValue) iter.Seq2[Pointer, JsonValue] {
	return func(yield func(Pointer, JsonValue) bool) {
		v.Accept(DfsVisitor(&ValueVisitor{yield: yield}))
	}
}
func (v *RootValue) DepthFirst() iter.Seq2[Pointer, JsonValue] { return depthFirstValues(v) }
func (v *Object) DepthFirst() iter.Seq2[Pointer, JsonValue]    { return depthFirstValues(v) }
func (v *Array) DepthFirst() iter.Seq2[Pointer, JsonValue]     { return depthFirstValues(v) }
func (v *String) DepthFirst() iter.Seq2[Pointer, JsonValue]    { return depthFirstValues(v) }
func (v *Number) DepthFirst() iter.Seq2[Pointer, JsonValue]    { return depthFirstValues(v) }
func (v *Bool) DepthFirst() iter.Seq2[Pointer, JsonValue]      { return depthFirstValues(v) }
func (v *Null) DepthFirst() iter.Seq2[Pointer, JsonValue]      { return depthFirstValues(v) }

func (vv *ValueVisitor) VisitObject(o *Object) error {
	vv.yield(vv.GetPointer(), o)
	return nil
}
func (vv *ValueVisitor) VisitArray(a *Array) error {
	vv.yield(vv.GetPointer(), a)
	return nil
}
func (vv *ValueVisitor) VisitString(s *String) error {
	vv.yield(vv.GetPointer(), s)
	return nil
}
func (vv *ValueVisitor) VisitNumber(n *Number) error {
	vv.yield(vv.GetPointer(), n)
	return nil
}
func (vv *ValueVisitor) VisitBool(b *Bool) error {
	vv.yield(vv.GetPointer(), b)
	return nil
}
func (vv *ValueVisitor) VisitNull(n *Null) error {
	vv.yield(vv.GetPointer(), n)
	return nil
}
