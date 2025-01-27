package fluffyjson

type (
	AccessAs interface {
		// TODO use pointer ?
		AccessAsObject
		AccessAsArray
		AccessAsString
		AccessAsNumber
		AccessAsBool
		AccessAsNull
	}
	AccessAsObject interface {
		AccessAsObject(Accessor) (Object, error)
		SliceAsObject(SliceAccessor) ([]Object, error)
		PointerAsObject(...Accessor) (Object, error)
	}
	AccessAsArray interface {
		AccessAsArray(Accessor) (Array, error)
		SliceAsArray(SliceAccessor) ([]Array, error)
		PointerAsArray(...Accessor) (Array, error)
	}

	AccessAsString interface {
		AccessAsString(Accessor) (String, error)
		SliceAsString(SliceAccessor) ([]String, error)
		PointerAsString(...Accessor) (String, error)
	}

	AccessAsNumber interface {
		AccessAsNumber(Accessor) (Number, error)
		SliceAsNumber(SliceAccessor) ([]Number, error)
		PointerAsNumber(...Accessor) (Number, error)
	}

	AccessAsBool interface {
		AccessAsBool(Accessor) (Bool, error)
		SliceAsBool(SliceAccessor) ([]Bool, error)
		PointerAsBool(...Accessor) (Bool, error)
	}

	AccessAsNull interface {
		AccessAsNull(Accessor) (Null, error)
		SliceAsNull(SliceAccessor) ([]Null, error)
		PointerAsNull(...Accessor) (Null, error)
	}
)

func accessAsObject(v JsonValue, acc Accessor) (Object, error) {
	v, err := acc.Accessing(v)
	if err != nil {
		return nil, err
	}
	o, err := v.AsObject()
	if err != nil {
		return nil, err
	}
	return o, nil
}
func accessAsArray(v JsonValue, acc Accessor) (Array, error) {
	v, err := acc.Accessing(v)
	if err != nil {
		return nil, err
	}
	arr, err := v.AsArray()
	if err != nil {
		return nil, err
	}
	return arr, nil
}
func accessAsString(v JsonValue, acc Accessor) (String, error) {
	v, err := acc.Accessing(v)
	if err != nil {
		return "", err
	}
	s, err := v.AsString()
	if err != nil {
		return "", err
	}
	return s, nil
}
func accessAsNumber(v JsonValue, acc Accessor) (Number, error) {
	v, err := acc.Accessing(v)
	if err != nil {
		return 0, err
	}
	n, err := v.AsNumber()
	if err != nil {
		return 0, err
	}
	return n, nil
}
func accessAsBool(v JsonValue, acc Accessor) (Bool, error) {
	v, err := acc.Accessing(v)
	if err != nil {
		return false, err
	}
	b, err := v.AsBool()
	if err != nil {
		return false, err
	}
	return b, nil
}
func accessAsNull(v JsonValue, acc Accessor) (Null, error) {
	v, err := acc.Accessing(v)
	if err != nil {
		return nil, err
	}
	n, err := v.AsNull()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (o Object) AccessAsObject(acc Accessor) (Object, error) { return accessAsObject(&o, acc) }
func (o Object) AccessAsArray(acc Accessor) (Array, error)   { return accessAsArray(&o, acc) }
func (o Object) AccessAsString(acc Accessor) (String, error) { return accessAsString(&o, acc) }
func (o Object) AccessAsNumber(acc Accessor) (Number, error) { return accessAsNumber(&o, acc) }
func (o Object) AccessAsBool(acc Accessor) (Bool, error)     { return accessAsBool(&o, acc) }
func (o Object) AccessAsNull(acc Accessor) (Null, error)     { return accessAsNull(&o, acc) }
func (a Array) AccessAsObject(acc Accessor) (Object, error)  { return accessAsObject(&a, acc) }
func (a Array) AccessAsArray(acc Accessor) (Array, error)    { return accessAsArray(&a, acc) }
func (a Array) AccessAsString(acc Accessor) (String, error)  { return accessAsString(&a, acc) }
func (a Array) AccessAsNumber(acc Accessor) (Number, error)  { return accessAsNumber(&a, acc) }
func (a Array) AccessAsBool(acc Accessor) (Bool, error)      { return accessAsBool(&a, acc) }
func (a Array) AccessAsNull(acc Accessor) (Null, error)      { return accessAsNull(&a, acc) }
func (s String) AccessAsObject(acc Accessor) (Object, error) { return accessAsObject(&s, acc) }
func (s String) AccessAsArray(acc Accessor) (Array, error)   { return accessAsArray(&s, acc) }
func (s String) AccessAsString(acc Accessor) (String, error) { return accessAsString(&s, acc) }
func (s String) AccessAsNumber(acc Accessor) (Number, error) { return accessAsNumber(&s, acc) }
func (s String) AccessAsBool(acc Accessor) (Bool, error)     { return accessAsBool(&s, acc) }
func (s String) AccessAsNull(acc Accessor) (Null, error)     { return accessAsNull(&s, acc) }
func (n Number) AccessAsObject(acc Accessor) (Object, error) { return accessAsObject(&n, acc) }
func (n Number) AccessAsArray(acc Accessor) (Array, error)   { return accessAsArray(&n, acc) }
func (n Number) AccessAsString(acc Accessor) (String, error) { return accessAsString(&n, acc) }
func (n Number) AccessAsNumber(acc Accessor) (Number, error) { return accessAsNumber(&n, acc) }
func (n Number) AccessAsBool(acc Accessor) (Bool, error)     { return accessAsBool(&n, acc) }
func (n Number) AccessAsNull(acc Accessor) (Null, error)     { return accessAsNull(&n, acc) }
func (b Bool) AccessAsObject(acc Accessor) (Object, error)   { return accessAsObject(&b, acc) }
func (b Bool) AccessAsArray(acc Accessor) (Array, error)     { return accessAsArray(&b, acc) }
func (b Bool) AccessAsString(acc Accessor) (String, error)   { return accessAsString(&b, acc) }
func (b Bool) AccessAsNumber(acc Accessor) (Number, error)   { return accessAsNumber(&b, acc) }
func (b Bool) AccessAsBool(acc Accessor) (Bool, error)       { return accessAsBool(&b, acc) }
func (b Bool) AccessAsNull(acc Accessor) (Null, error)       { return accessAsNull(&b, acc) }
func (n Null) AccessAsObject(acc Accessor) (Object, error)   { return accessAsObject(&n, acc) }
func (n Null) AccessAsArray(acc Accessor) (Array, error)     { return accessAsArray(&n, acc) }
func (n Null) AccessAsString(acc Accessor) (String, error)   { return accessAsString(&n, acc) }
func (n Null) AccessAsNumber(acc Accessor) (Number, error)   { return accessAsNumber(&n, acc) }
func (n Null) AccessAsBool(acc Accessor) (Bool, error)       { return accessAsBool(&n, acc) }
func (n Null) AccessAsNull(acc Accessor) (Null, error)       { return accessAsNull(&n, acc) }

func sliceAsObject(v JsonValue, acc SliceAccessor) ([]Object, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	s := make([]Object, 0, len(vs))
	for _, v := range vs {
		o, err := v.AsObject()
		if err != nil {
			return nil, err
		}
		s = append(s, o)
	}
	return s, nil
}
func sliceAsArray(v JsonValue, acc SliceAccessor) ([]Array, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	s := make([]Array, 0, len(vs))
	for _, v := range vs {
		a, err := v.AsArray()
		if err != nil {
			return nil, err
		}
		s = append(s, a)
	}
	return s, nil
}
func sliceAsString(v JsonValue, acc SliceAccessor) ([]String, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	s := make([]String, 0, len(vs))
	for _, v := range vs {
		a, err := v.AsString()
		if err != nil {
			return nil, err
		}
		s = append(s, a)
	}
	return s, nil
}
func sliceAsNumber(v JsonValue, acc SliceAccessor) ([]Number, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	s := make([]Number, 0, len(vs))
	for _, v := range vs {
		a, err := v.AsNumber()
		if err != nil {
			return nil, err
		}
		s = append(s, a)
	}
	return s, nil
}
func sliceAsBool(v JsonValue, acc SliceAccessor) ([]Bool, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	s := make([]Bool, 0, len(vs))
	for _, v := range vs {
		a, err := v.AsBool()
		if err != nil {
			return nil, err
		}
		s = append(s, a)
	}
	return s, nil
}
func sliceAsNull(v JsonValue, acc SliceAccessor) ([]Null, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	s := make([]Null, 0, len(vs))
	for _, v := range vs {
		a, err := v.AsNull()
		if err != nil {
			return nil, err
		}
		s = append(s, a)
	}
	return s, nil
}
func (o Object) SliceAsObject(acc SliceAccessor) ([]Object, error) { return sliceAsObject(&o, acc) }
func (o Object) SliceAsArray(acc SliceAccessor) ([]Array, error)   { return sliceAsArray(&o, acc) }
func (o Object) SliceAsString(acc SliceAccessor) ([]String, error) { return sliceAsString(&o, acc) }
func (o Object) SliceAsNumber(acc SliceAccessor) ([]Number, error) { return sliceAsNumber(&o, acc) }
func (o Object) SliceAsBool(acc SliceAccessor) ([]Bool, error)     { return sliceAsBool(&o, acc) }
func (o Object) SliceAsNull(acc SliceAccessor) ([]Null, error)     { return sliceAsNull(&o, acc) }
func (a Array) SliceAsObject(acc SliceAccessor) ([]Object, error)  { return sliceAsObject(&a, acc) }
func (a Array) SliceAsArray(acc SliceAccessor) ([]Array, error)    { return sliceAsArray(&a, acc) }
func (a Array) SliceAsString(acc SliceAccessor) ([]String, error)  { return sliceAsString(&a, acc) }
func (a Array) SliceAsNumber(acc SliceAccessor) ([]Number, error)  { return sliceAsNumber(&a, acc) }
func (a Array) SliceAsBool(acc SliceAccessor) ([]Bool, error)      { return sliceAsBool(&a, acc) }
func (a Array) SliceAsNull(acc SliceAccessor) ([]Null, error)      { return sliceAsNull(&a, acc) }
func (s String) SliceAsObject(acc SliceAccessor) ([]Object, error) { return sliceAsObject(&s, acc) }
func (s String) SliceAsArray(acc SliceAccessor) ([]Array, error)   { return sliceAsArray(&s, acc) }
func (s String) SliceAsString(acc SliceAccessor) ([]String, error) { return sliceAsString(&s, acc) }
func (s String) SliceAsNumber(acc SliceAccessor) ([]Number, error) { return sliceAsNumber(&s, acc) }
func (s String) SliceAsBool(acc SliceAccessor) ([]Bool, error)     { return sliceAsBool(&s, acc) }
func (s String) SliceAsNull(acc SliceAccessor) ([]Null, error)     { return sliceAsNull(&s, acc) }
func (n Number) SliceAsObject(acc SliceAccessor) ([]Object, error) { return sliceAsObject(&n, acc) }
func (n Number) SliceAsArray(acc SliceAccessor) ([]Array, error)   { return sliceAsArray(&n, acc) }
func (n Number) SliceAsString(acc SliceAccessor) ([]String, error) { return sliceAsString(&n, acc) }
func (n Number) SliceAsNumber(acc SliceAccessor) ([]Number, error) { return sliceAsNumber(&n, acc) }
func (n Number) SliceAsBool(acc SliceAccessor) ([]Bool, error)     { return sliceAsBool(&n, acc) }
func (n Number) SliceAsNull(acc SliceAccessor) ([]Null, error)     { return sliceAsNull(&n, acc) }
func (b Bool) SliceAsObject(acc SliceAccessor) ([]Object, error)   { return sliceAsObject(&b, acc) }
func (b Bool) SliceAsArray(acc SliceAccessor) ([]Array, error)     { return sliceAsArray(&b, acc) }
func (b Bool) SliceAsString(acc SliceAccessor) ([]String, error)   { return sliceAsString(&b, acc) }
func (b Bool) SliceAsNumber(acc SliceAccessor) ([]Number, error)   { return sliceAsNumber(&b, acc) }
func (b Bool) SliceAsBool(acc SliceAccessor) ([]Bool, error)       { return sliceAsBool(&b, acc) }
func (b Bool) SliceAsNull(acc SliceAccessor) ([]Null, error)       { return sliceAsNull(&b, acc) }
func (n Null) SliceAsObject(acc SliceAccessor) ([]Object, error)   { return sliceAsObject(&n, acc) }
func (n Null) SliceAsArray(acc SliceAccessor) ([]Array, error)     { return sliceAsArray(&n, acc) }
func (n Null) SliceAsString(acc SliceAccessor) ([]String, error)   { return sliceAsString(&n, acc) }
func (n Null) SliceAsNumber(acc SliceAccessor) ([]Number, error)   { return sliceAsNumber(&n, acc) }
func (n Null) SliceAsBool(acc SliceAccessor) ([]Bool, error)       { return sliceAsBool(&n, acc) }
func (n Null) SliceAsNull(acc SliceAccessor) ([]Null, error)       { return sliceAsNull(&n, acc) }

func pointerAsObject(v JsonValue, ptr Pointer) (Object, error) {
	v, err := ptr.Accessing(v)
	if err != nil {
		return nil, err
	}
	o, err := v.AsObject()
	if err != nil {
		return nil, err
	}
	return o, nil
}
func pointerAsArray(v JsonValue, ptr Pointer) (Array, error) {
	v, err := ptr.Accessing(v)
	if err != nil {
		return nil, err
	}
	a, err := v.AsArray()
	if err != nil {
		return nil, err
	}
	return a, nil
}
func pointerAsString(v JsonValue, ptr Pointer) (String, error) {
	v, err := ptr.Accessing(v)
	if err != nil {
		return "", err
	}
	s, err := v.AsString()
	if err != nil {
		return "", err
	}
	return s, nil
}
func pointerAsNumber(v JsonValue, ptr Pointer) (Number, error) {
	v, err := ptr.Accessing(v)
	if err != nil {
		return 0, err
	}
	n, err := v.AsNumber()
	if err != nil {
		return 0, err
	}
	return n, nil
}
func pointerAsBool(v JsonValue, ptr Pointer) (Bool, error) {
	v, err := ptr.Accessing(v)
	if err != nil {
		return false, err
	}
	b, err := v.AsBool()
	if err != nil {
		return false, err
	}
	return b, nil
}
func pointerAsNull(v JsonValue, ptr Pointer) (Null, error) {
	v, err := ptr.Accessing(v)
	if err != nil {
		return nil, err
	}
	n, err := v.AsNull()
	if err != nil {
		return nil, err
	}
	return n, nil
}
func (o Object) PointerAsObject(ptr ...Accessor) (Object, error) { return pointerAsObject(&o, ptr) }
func (o Object) PointerAsArray(ptr ...Accessor) (Array, error)   { return pointerAsArray(&o, ptr) }
func (o Object) PointerAsString(ptr ...Accessor) (String, error) { return pointerAsString(&o, ptr) }
func (o Object) PointerAsNumber(ptr ...Accessor) (Number, error) { return pointerAsNumber(&o, ptr) }
func (o Object) PointerAsBool(ptr ...Accessor) (Bool, error)     { return pointerAsBool(&o, ptr) }
func (o Object) PointerAsNull(ptr ...Accessor) (Null, error)     { return pointerAsNull(&o, ptr) }
func (a Array) PointerAsObject(ptr ...Accessor) (Object, error)  { return pointerAsObject(&a, ptr) }
func (a Array) PointerAsArray(ptr ...Accessor) (Array, error)    { return pointerAsArray(&a, ptr) }
func (a Array) PointerAsString(ptr ...Accessor) (String, error)  { return pointerAsString(&a, ptr) }
func (a Array) PointerAsNumber(ptr ...Accessor) (Number, error)  { return pointerAsNumber(&a, ptr) }
func (a Array) PointerAsBool(ptr ...Accessor) (Bool, error)      { return pointerAsBool(&a, ptr) }
func (a Array) PointerAsNull(ptr ...Accessor) (Null, error)      { return pointerAsNull(&a, ptr) }
func (s String) PointerAsObject(ptr ...Accessor) (Object, error) { return pointerAsObject(&s, ptr) }
func (s String) PointerAsArray(ptr ...Accessor) (Array, error)   { return pointerAsArray(&s, ptr) }
func (s String) PointerAsString(ptr ...Accessor) (String, error) { return pointerAsString(&s, ptr) }
func (s String) PointerAsNumber(ptr ...Accessor) (Number, error) { return pointerAsNumber(&s, ptr) }
func (s String) PointerAsBool(ptr ...Accessor) (Bool, error)     { return pointerAsBool(&s, ptr) }
func (s String) PointerAsNull(ptr ...Accessor) (Null, error)     { return pointerAsNull(&s, ptr) }
func (n Number) PointerAsObject(ptr ...Accessor) (Object, error) { return pointerAsObject(&n, ptr) }
func (n Number) PointerAsArray(ptr ...Accessor) (Array, error)   { return pointerAsArray(&n, ptr) }
func (n Number) PointerAsString(ptr ...Accessor) (String, error) { return pointerAsString(&n, ptr) }
func (n Number) PointerAsNumber(ptr ...Accessor) (Number, error) { return pointerAsNumber(&n, ptr) }
func (n Number) PointerAsBool(ptr ...Accessor) (Bool, error)     { return pointerAsBool(&n, ptr) }
func (n Number) PointerAsNull(ptr ...Accessor) (Null, error)     { return pointerAsNull(&n, ptr) }
func (b Bool) PointerAsObject(ptr ...Accessor) (Object, error)   { return pointerAsObject(&b, ptr) }
func (b Bool) PointerAsArray(ptr ...Accessor) (Array, error)     { return pointerAsArray(&b, ptr) }
func (b Bool) PointerAsString(ptr ...Accessor) (String, error)   { return pointerAsString(&b, ptr) }
func (b Bool) PointerAsNumber(ptr ...Accessor) (Number, error)   { return pointerAsNumber(&b, ptr) }
func (b Bool) PointerAsBool(ptr ...Accessor) (Bool, error)       { return pointerAsBool(&b, ptr) }
func (b Bool) PointerAsNull(ptr ...Accessor) (Null, error)       { return pointerAsNull(&b, ptr) }
func (n Null) PointerAsObject(ptr ...Accessor) (Object, error)   { return pointerAsObject(&n, ptr) }
func (n Null) PointerAsArray(ptr ...Accessor) (Array, error)     { return pointerAsArray(&n, ptr) }
func (n Null) PointerAsString(ptr ...Accessor) (String, error)   { return pointerAsString(&n, ptr) }
func (n Null) PointerAsNumber(ptr ...Accessor) (Number, error)   { return pointerAsNumber(&n, ptr) }
func (n Null) PointerAsBool(ptr ...Accessor) (Bool, error)       { return pointerAsBool(&n, ptr) }
func (n Null) PointerAsNull(ptr ...Accessor) (Null, error)       { return pointerAsNull(&n, ptr) }
