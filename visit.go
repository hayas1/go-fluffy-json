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
		BreadthFirst() iter.Seq2[Pointer, JsonValue]
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
	Bfs[V Visitor] struct {
		pointerBuf []Pointer
		valueBuf   []JsonValue
		visitor    V
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
	defer func() { err = dfs.LeaveRoot(v) }() // TODO if err != nil, leave should be called or not ?
	if err := v.JsonValue.Accept(dfs); err != nil {
		return err
	}
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
	dfs.SetPointer(append(dfs.GetPointer(), KeyAccess(k)))
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
	defer dfs.SetPointer(dfs.GetPointer()[:len(dfs.GetPointer())-1])
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
	dfs.SetPointer(append(dfs.GetPointer(), IndexAccess(i)))
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
	defer dfs.SetPointer(dfs.GetPointer()[:len(dfs.GetPointer())-1])
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

func BfsVisitor[V Visitor](visitor V) *Bfs[V] {
	return &Bfs[V]{visitor: visitor}
}
func (bfs *Bfs[V]) GetPointer() Pointer  { return bfs.visitor.GetPointer() }
func (bfs *Bfs[V]) SetPointer(p Pointer) { bfs.visitor.SetPointer(p) }
func (bfs *Bfs[V]) VisitRoot(v *RootValue) (err error) {
	if err = bfs.visitor.VisitRoot(v); err != nil {
		return err
	}
	defer func() { err = bfs.LeaveRoot(v) }() // TODO if err != nil, leave should be called or not ?

	bfs.pointerBuf, bfs.valueBuf = append(bfs.pointerBuf, nil), append(bfs.valueBuf, v.JsonValue)
	for len(bfs.pointerBuf) > 0 {
		p, v := bfs.pointerBuf[0], bfs.valueBuf[0]
		bfs.pointerBuf, bfs.valueBuf = bfs.pointerBuf[1:], bfs.valueBuf[1:]

		bfs.visitor.SetPointer(p)
		if err := v.Accept(bfs); err != nil {
			return err
		}
	}
	return nil
}
func (bfs *Bfs[V]) LeaveRoot(v *RootValue) error {
	return bfs.visitor.LeaveRoot(v)
}
func (bfs *Bfs[V]) VisitObject(o *Object) (err error) {
	if err = bfs.visitor.VisitObject(o); err != nil {
		return err
	}
	defer func() { err = bfs.LeaveObject(o) }()

	for k, v := range *o {
		bfs.pointerBuf = append(bfs.pointerBuf, append(bfs.visitor.GetPointer(), KeyAccess(k)))
		bfs.valueBuf = append(bfs.valueBuf, v)
	}
	return nil
}
func (bfs *Bfs[V]) LeaveObject(o *Object) error {
	return bfs.visitor.LeaveObject(o)
}
func (bfs *Bfs[V]) VisitObjectEntry(k string, v JsonValue) (err error) {
	if err = bfs.visitor.VisitObjectEntry(k, v); err != nil {
		return err
	}
	defer func() { err = bfs.LeaveObjectEntry(k, v) }()

	return nil
}
func (bfs *Bfs[V]) LeaveObjectEntry(k string, v JsonValue) error {
	return bfs.visitor.LeaveObjectEntry(k, v)
}
func (bfs *Bfs[V]) VisitArray(a *Array) (err error) {
	if err = bfs.visitor.VisitArray(a); err != nil {
		return err
	}
	defer func() { err = bfs.LeaveArray(a) }()

	for i, v := range *a {
		bfs.pointerBuf = append(bfs.pointerBuf, append(bfs.visitor.GetPointer(), IndexAccess(i)))
		bfs.valueBuf = append(bfs.valueBuf, v)
	}
	return nil
}
func (bfs *Bfs[V]) LeaveArray(a *Array) error {
	return bfs.visitor.LeaveArray(a)
}
func (bfs *Bfs[V]) VisitArrayEntry(i int, v JsonValue) (err error) {
	if err = bfs.visitor.VisitArrayEntry(i, v); err != nil {
		return err
	}
	defer func() { err = bfs.LeaveArrayEntry(i, v) }()

	return nil
}
func (bfs *Bfs[V]) LeaveArrayEntry(i int, v JsonValue) error {
	return bfs.visitor.LeaveArrayEntry(i, v)
}
func (bfs *Bfs[V]) VisitString(s *String) error {
	return bfs.visitor.VisitString(s)
}
func (bfs *Bfs[V]) VisitNumber(n *Number) error {
	return bfs.visitor.VisitNumber(n)
}
func (bfs *Bfs[V]) VisitBool(b *Bool) error {
	return bfs.visitor.VisitBool(b)
}
func (bfs *Bfs[V]) VisitNull(n *Null) error {
	return bfs.visitor.VisitNull(n)
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

func breadthFirstValues(v JsonValue) iter.Seq2[Pointer, JsonValue] {
	return func(yield func(Pointer, JsonValue) bool) {
		v.Accept(BfsVisitor(&ValueVisitor{yield: yield}))
	}
}
func (v *RootValue) BreadthFirst() iter.Seq2[Pointer, JsonValue] { return breadthFirstValues(v) }
func (v *Object) BreadthFirst() iter.Seq2[Pointer, JsonValue]    { return breadthFirstValues(v) }
func (v *Array) BreadthFirst() iter.Seq2[Pointer, JsonValue]     { return breadthFirstValues(v) }
func (v *String) BreadthFirst() iter.Seq2[Pointer, JsonValue]    { return breadthFirstValues(v) }
func (v *Number) BreadthFirst() iter.Seq2[Pointer, JsonValue]    { return breadthFirstValues(v) }
func (v *Bool) BreadthFirst() iter.Seq2[Pointer, JsonValue]      { return breadthFirstValues(v) }
func (v *Null) BreadthFirst() iter.Seq2[Pointer, JsonValue]      { return breadthFirstValues(v) }

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
