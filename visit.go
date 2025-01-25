package fluffyjson

type (
	Accept interface {
		Accept(Visitor) error
	}
	Visitor interface {
		VisitObject(*Object) error
		VisitObjectEntry(string, JsonValue) error
		LeaveObjectEntry(string, JsonValue) error
		LeaveObject(*Object) error
		VisitArray(*Array) error
		VisitArrayEntry(int, JsonValue) error
		LeaveArrayEntry(int, JsonValue) error
		LeaveArray(*Array) error
		VisitString(*String) error
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

func (v *BaseVisitor) VisitObject(o *Object) error                      { return nil }
func (v *BaseVisitor) VisitObjectEntry(key string, val JsonValue) error { return nil }
func (v *BaseVisitor) LeaveObjectEntry(key string, val JsonValue) error { return nil }
func (v *BaseVisitor) LeaveObject(o *Object) error                      { return nil }
func (v *BaseVisitor) VisitArray(a *Array) error                        { return nil }
func (v *BaseVisitor) VisitArrayEntry(idx int, val JsonValue) error     { return nil }
func (v *BaseVisitor) LeaveArrayEntry(idx int, val JsonValue) error     { return nil }
func (v *BaseVisitor) LeaveArray(a *Array) error                        { return nil }
func (v *BaseVisitor) VisitString(s *String) error                      { return nil }

func DfsVisitor[V Visitor](visitor V) *Dfs[V] {
	return &Dfs[V]{Visitor: visitor}
}
func (dfs *Dfs[V]) VisitObject(o *Object) error {
	dfs.Visitor.VisitObject(o)
	for k, v := range *o {
		if err := dfs.VisitObjectEntry(k, v); err != nil {
			return err
		}
	}
	dfs.Visitor.LeaveObject(o)
	return nil
}
func (dfs *Dfs[V]) VisitObjectEntry(k string, v JsonValue) error {
	dfs.Visitor.VisitObjectEntry(k, v)
	if err := v.Accept(dfs); err != nil {
		return err
	}
	dfs.Visitor.LeaveObjectEntry(k, v)
	return nil
}
func (dfs *Dfs[V]) LeaveObjectEntry(k string, v JsonValue) error {
	return nil
}
func (dfs *Dfs[V]) LeaveObject(o *Object) error {
	return nil
}
func (dfs *Dfs[V]) VisitArray(a *Array) error {
	dfs.Visitor.VisitArray(a)
	for i, v := range *a {
		if err := dfs.VisitArrayEntry(i, v); err != nil {
			return err
		}
	}
	dfs.Visitor.LeaveArray(a)
	return nil
}
func (dfs *Dfs[V]) VisitArrayEntry(i int, v JsonValue) error {
	dfs.Visitor.VisitArrayEntry(i, v)
	if err := v.Accept(dfs); err != nil {
		return err
	}
	dfs.Visitor.LeaveArrayEntry(i, v)
	return nil
}
func (dfs *Dfs[V]) LeaveArrayEntry(i int, v JsonValue) error {
	return nil
}
func (dfs *Dfs[V]) LeaveArray(a *Array) error {
	return nil
}
func (dfs *Dfs[V]) VisitString(s *String) error {
	return dfs.Visitor.VisitString(s)
}
