package fluffyjson

type (
	Accept interface {
		Accept(Visitor) error
	}
	Visitor interface {
		VisitObject(*Object) error
		LeaveObject(*Object) error
		VisitArray(*Array) error
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

func (v *BaseVisitor) VisitObject(o *Object) error { return nil }
func (v *BaseVisitor) LeaveObject(o *Object) error { return nil }
func (v *BaseVisitor) VisitArray(a *Array) error   { return nil }
func (v *BaseVisitor) LeaveArray(a *Array) error   { return nil }
func (v *BaseVisitor) VisitString(s *String) error { return nil }

func DfsVisitor[V Visitor](visitor V) *Dfs[V] {
	return &Dfs[V]{Visitor: visitor}
}
func (dfs *Dfs[V]) VisitObject(o *Object) error {
	dfs.Visitor.VisitObject(o)
	for _, v := range *o {
		if err := v.Accept(dfs); err != nil {
			return err
		}
	}
	dfs.Visitor.LeaveObject(o)
	return nil
}
func (dfs *Dfs[V]) LeaveObject(o *Object) error {
	return nil
}
func (dfs *Dfs[V]) VisitArray(a *Array) error {
	dfs.Visitor.VisitArray(a)
	for _, v := range *a {
		if err := v.Accept(dfs); err != nil {
			return err
		}
	}
	dfs.Visitor.LeaveArray(a)
	return nil
}
func (dfs *Dfs[V]) LeaveArray(a *Array) error {
	return nil
}
func (dfs *Dfs[V]) VisitString(s *String) error {
	return dfs.Visitor.VisitString(s)
}
