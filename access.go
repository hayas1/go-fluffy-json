package fluffyjson

type (
	AccessAs interface {
		Access
		AccessAsObject
		AccessAsArray
		AccessAsString
		AccessAsNumber
		AccessAsBool
		AccessAsNull
	}
	AccessAsObject interface {
		AccessAsObject(...Accessor) (Object, error)
		SliceAsObject(SliceAccessor) ([]Object, error)
	}
	AccessAsArray interface {
		AccessAsArray(...Accessor) (Array, error)
		SliceAsArray(SliceAccessor) ([]Array, error)
	}

	AccessAsString interface {
		AccessAsString(...Accessor) (String, error)
		SliceAsString(SliceAccessor) ([]String, error)
	}

	AccessAsNumber interface {
		AccessAsNumber(...Accessor) (Number, error)
		SliceAsNumber(SliceAccessor) ([]Number, error)
	}

	AccessAsBool interface {
		AccessAsBool(...Accessor) (Bool, error)
		SliceAsBool(SliceAccessor) ([]Bool, error)
	}

	AccessAsNull interface {
		AccessAsNull(...Accessor) (Null, error)
		SliceAsNull(SliceAccessor) ([]Null, error)
	}
)

func accessAsObject(v JsonValue, ptr ...Accessor) (Object, error) {
	v, err := Pointer(ptr).Accessing(v)
	if err != nil {
		return nil, err
	}
	return v.AsObject()
}
func accessAsArray(v JsonValue, ptr ...Accessor) (Array, error) {
	v, err := Pointer(ptr).Accessing(v)
	if err != nil {
		return nil, err
	}
	return v.AsArray()
}
func accessAsString(v JsonValue, ptr ...Accessor) (String, error) {
	v, err := Pointer(ptr).Accessing(v)
	if err != nil {
		return "", err
	}
	return v.AsString()
}
func accessAsNumber(v JsonValue, ptr ...Accessor) (Number, error) {
	v, err := Pointer(ptr).Accessing(v)
	if err != nil {
		return 0, err
	}
	return v.AsNumber()
}
func accessAsBool(v JsonValue, ptr ...Accessor) (Bool, error) {
	v, err := Pointer(ptr).Accessing(v)
	if err != nil {
		return false, err
	}
	return v.AsBool()
}
func accessAsNull(v JsonValue, ptr ...Accessor) (Null, error) {
	v, err := Pointer(ptr).Accessing(v)
	if err != nil {
		return nil, err
	}
	return v.AsNull()
}

func (o Object) AccessAsObject(ptr ...Accessor) (Object, error) { return accessAsObject(&o, ptr...) }
func (o Object) AccessAsArray(ptr ...Accessor) (Array, error)   { return accessAsArray(&o, ptr...) }
func (o Object) AccessAsString(ptr ...Accessor) (String, error) { return accessAsString(&o, ptr...) }
func (o Object) AccessAsNumber(ptr ...Accessor) (Number, error) { return accessAsNumber(&o, ptr...) }
func (o Object) AccessAsBool(ptr ...Accessor) (Bool, error)     { return accessAsBool(&o, ptr...) }
func (o Object) AccessAsNull(ptr ...Accessor) (Null, error)     { return accessAsNull(&o, ptr...) }
func (a Array) AccessAsObject(ptr ...Accessor) (Object, error)  { return accessAsObject(&a, ptr...) }
func (a Array) AccessAsArray(ptr ...Accessor) (Array, error)    { return accessAsArray(&a, ptr...) }
func (a Array) AccessAsString(ptr ...Accessor) (String, error)  { return accessAsString(&a, ptr...) }
func (a Array) AccessAsNumber(ptr ...Accessor) (Number, error)  { return accessAsNumber(&a, ptr...) }
func (a Array) AccessAsBool(ptr ...Accessor) (Bool, error)      { return accessAsBool(&a, ptr...) }
func (a Array) AccessAsNull(ptr ...Accessor) (Null, error)      { return accessAsNull(&a, ptr...) }
func (s String) AccessAsObject(ptr ...Accessor) (Object, error) { return accessAsObject(&s, ptr...) }
func (s String) AccessAsArray(ptr ...Accessor) (Array, error)   { return accessAsArray(&s, ptr...) }
func (s String) AccessAsString(ptr ...Accessor) (String, error) { return accessAsString(&s, ptr...) }
func (s String) AccessAsNumber(ptr ...Accessor) (Number, error) { return accessAsNumber(&s, ptr...) }
func (s String) AccessAsBool(ptr ...Accessor) (Bool, error)     { return accessAsBool(&s, ptr...) }
func (s String) AccessAsNull(ptr ...Accessor) (Null, error)     { return accessAsNull(&s, ptr...) }
func (n Number) AccessAsObject(ptr ...Accessor) (Object, error) { return accessAsObject(&n, ptr...) }
func (n Number) AccessAsArray(ptr ...Accessor) (Array, error)   { return accessAsArray(&n, ptr...) }
func (n Number) AccessAsString(ptr ...Accessor) (String, error) { return accessAsString(&n, ptr...) }
func (n Number) AccessAsNumber(ptr ...Accessor) (Number, error) { return accessAsNumber(&n, ptr...) }
func (n Number) AccessAsBool(ptr ...Accessor) (Bool, error)     { return accessAsBool(&n, ptr...) }
func (n Number) AccessAsNull(ptr ...Accessor) (Null, error)     { return accessAsNull(&n, ptr...) }
func (b Bool) AccessAsObject(ptr ...Accessor) (Object, error)   { return accessAsObject(&b, ptr...) }
func (b Bool) AccessAsArray(ptr ...Accessor) (Array, error)     { return accessAsArray(&b, ptr...) }
func (b Bool) AccessAsString(ptr ...Accessor) (String, error)   { return accessAsString(&b, ptr...) }
func (b Bool) AccessAsNumber(ptr ...Accessor) (Number, error)   { return accessAsNumber(&b, ptr...) }
func (b Bool) AccessAsBool(ptr ...Accessor) (Bool, error)       { return accessAsBool(&b, ptr...) }
func (b Bool) AccessAsNull(ptr ...Accessor) (Null, error)       { return accessAsNull(&b, ptr...) }
func (n Null) AccessAsObject(ptr ...Accessor) (Object, error)   { return accessAsObject(&n, ptr...) }
func (n Null) AccessAsArray(ptr ...Accessor) (Array, error)     { return accessAsArray(&n, ptr...) }
func (n Null) AccessAsString(ptr ...Accessor) (String, error)   { return accessAsString(&n, ptr...) }
func (n Null) AccessAsNumber(ptr ...Accessor) (Number, error)   { return accessAsNumber(&n, ptr...) }
func (n Null) AccessAsBool(ptr ...Accessor) (Bool, error)       { return accessAsBool(&n, ptr...) }
func (n Null) AccessAsNull(ptr ...Accessor) (Null, error)       { return accessAsNull(&n, ptr...) }

func sliceAsObject(v JsonValue, acc SliceAccessor) ([]Object, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	slice := make([]Object, 0, len(vs))
	for _, v := range vs {
		o, err := v.AsObject()
		if err != nil {
			return nil, err
		}
		slice = append(slice, o)
	}
	return slice, nil
}
func sliceAsArray(v JsonValue, acc SliceAccessor) ([]Array, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	slice := make([]Array, 0, len(vs))
	for _, v := range vs {
		a, err := v.AsArray()
		if err != nil {
			return nil, err
		}
		slice = append(slice, a)
	}
	return slice, nil
}
func sliceAsString(v JsonValue, acc SliceAccessor) ([]String, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	slice := make([]String, 0, len(vs))
	for _, v := range vs {
		s, err := v.AsString()
		if err != nil {
			return nil, err
		}
		slice = append(slice, s)
	}
	return slice, nil
}
func sliceAsNumber(v JsonValue, acc SliceAccessor) ([]Number, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	slice := make([]Number, 0, len(vs))
	for _, v := range vs {
		n, err := v.AsNumber()
		if err != nil {
			return nil, err
		}
		slice = append(slice, n)
	}
	return slice, nil
}
func sliceAsBool(v JsonValue, acc SliceAccessor) ([]Bool, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	slice := make([]Bool, 0, len(vs))
	for _, v := range vs {
		b, err := v.AsBool()
		if err != nil {
			return nil, err
		}
		slice = append(slice, b)
	}
	return slice, nil
}
func sliceAsNull(v JsonValue, acc SliceAccessor) ([]Null, error) {
	vs, err := acc.Slicing(v)
	if err != nil {
		return nil, err
	}
	slice := make([]Null, 0, len(vs))
	for _, v := range vs {
		n, err := v.AsNull()
		if err != nil {
			return nil, err
		}
		slice = append(slice, n)
	}
	return slice, nil
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
