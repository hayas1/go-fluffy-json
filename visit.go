package fluffyjson

type (
	Accept interface {
		Accept(Visitor) error
	}
	Visitor interface {
		VisitObject(*Object) error
		// LeaveObject(*Object) error
		VisitArray(*Array) error
		// LeaveArray(*Array) error
		VisitString(*String) error
	}
	DfsVisitor struct {
		// Inner Visitor
	}
)

func (v *Value) Accept(visitor Visitor) error  { return v.Value.Accept(visitor) }
func (o *Object) Accept(visitor Visitor) error { return visitor.VisitObject(o) }
func (a *Array) Accept(visitor Visitor) error  { return visitor.VisitArray(a) }
func (s *String) Accept(visitor Visitor) error { return visitor.VisitString(s) }

func (dfs *DfsVisitor) VisitObject(o *Object) error {
	// o.Accept(dfs.Inner)
	for _, v := range *o {
		if err := v.Accept(dfs); err != nil {
			return err
		}
	}
	return nil
}
func (dfs *DfsVisitor) VisitArray(a *Array) error {
	// a.Accept(dfs.Inner)
	for _, v := range *a {
		if err := v.Accept(dfs); err != nil {
			return err
		}
	}
	return nil
}
func (dfs *DfsVisitor) VisitString(s *String) error {
	// return s.Accept(dfs.Inner)
	return nil
}
