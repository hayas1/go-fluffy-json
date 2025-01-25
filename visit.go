package fluffyjson

type (
	Accept interface {
		Accept(Visitor) error
	}
	Visitor interface {
		VisitObject(*Object) error
		VisitObjectEntry(*ObjectEntry) error
		LeaveObjectEntry(*ObjectEntry) error
		LeaveObject(*Object) error
		VisitArray(*Array) error
		VisitArrayEntry(*ArrayEntry) error
		LeaveArrayEntry(*ArrayEntry) error
		LeaveArray(*Array) error
		VisitString(*String) error
	}

	ObjectEntry struct {
		Key   string
		Value JsonValue
	}
	ArrayEntry struct {
		Index int
		Value JsonValue
	}

	BaseVisitor    struct{}
	Dfs[V Visitor] struct {
		Visitor V
	}
)

func (v *Value) Accept(visitor Visitor) error  { return v.Value.Accept(visitor) }
func (o *Object) Accept(visitor Visitor) error { visitor.VisitObject(o); return visitor.LeaveObject(o) }
func (a *Array) Accept(visitor Visitor) error  { visitor.VisitArray(a); return visitor.LeaveArray(a) }
func (s *String) Accept(visitor Visitor) error { return visitor.VisitString(s) }

func (v *BaseVisitor) VisitObject(o *Object) error           { return nil }
func (v *BaseVisitor) VisitObjectEntry(o *ObjectEntry) error { return nil }
func (v *BaseVisitor) LeaveObjectEntry(o *ObjectEntry) error { return nil }
func (v *BaseVisitor) LeaveObject(o *Object) error           { return nil }
func (v *BaseVisitor) VisitArray(a *Array) error             { return nil }
func (v *BaseVisitor) VisitArrayEntry(a *ArrayEntry) error   { return nil }
func (v *BaseVisitor) LeaveArrayEntry(a *ArrayEntry) error   { return nil }
func (v *BaseVisitor) LeaveArray(a *Array) error             { return nil }
func (v *BaseVisitor) VisitString(s *String) error           { return nil }

func DfsVisitor[V Visitor](visitor V) *Dfs[V] {
	return &Dfs[V]{Visitor: visitor}
}
func (dfs *Dfs[V]) VisitObject(o *Object) error {
	dfs.Visitor.VisitObject(o)
	for k, v := range *o {
		entry := ObjectEntry{Key: k, Value: v}
		if err := dfs.VisitObjectEntry(&entry); err != nil {
			return err
		}
	}
	dfs.Visitor.LeaveObject(o)
	return nil
}
func (dfs *Dfs[V]) VisitObjectEntry(e *ObjectEntry) error {
	dfs.Visitor.VisitObjectEntry(e)
	if err := e.Value.Accept(dfs); err != nil {
		return err
	}
	dfs.Visitor.LeaveObjectEntry(e)
	return nil
}
func (dfs *Dfs[V]) LeaveObjectEntry(e *ObjectEntry) error {
	return nil
}
func (dfs *Dfs[V]) LeaveObject(o *Object) error {
	return nil
}
func (dfs *Dfs[V]) VisitArray(a *Array) error {
	dfs.Visitor.VisitArray(a)
	for i, v := range *a {
		entry := ArrayEntry{Index: i, Value: v}
		if err := dfs.VisitArrayEntry(&entry); err != nil {
			return err
		}
	}
	dfs.Visitor.LeaveArray(a)
	return nil
}
func (dfs *Dfs[V]) VisitArrayEntry(a *ArrayEntry) error {
	dfs.Visitor.VisitArrayEntry(a)
	if err := a.Value.Accept(dfs); err != nil {
		return err
	}
	dfs.Visitor.LeaveArrayEntry(a)
	return nil
}
func (dfs *Dfs[V]) LeaveArrayEntry(a *ArrayEntry) error {
	return nil
}
func (dfs *Dfs[V]) LeaveArray(a *Array) error {
	return nil
}
func (dfs *Dfs[V]) VisitString(s *String) error {
	return dfs.Visitor.VisitString(s)
}
