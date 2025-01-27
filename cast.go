package fluffyjson

type (
	AsValue interface {
		AsObject
		AsArray
		AsString
		AsNumber
		AsBool
		AsNull
	}

	AsObject interface {
		IsObject() bool
		AsObject() (Object, error)
	}
	AsArray interface {
		IsArray() bool
		AsArray() (Array, error)
	}
	AsString interface {
		IsString() bool
		AsString() (String, error)
	}
	AsNumber interface {
		IsNumber() bool
		AsNumber() (Number, error)
	}
	AsBool interface {
		IsBool() bool
		AsBool() (Bool, error)
	}
	AsNull interface {
		IsNull() bool
		AsNull() (Null, error)
	}
)

func (o Object) IsObject() bool            { return true }
func (o Object) AsObject() (Object, error) { return o, nil }
func (o Object) IsArray() bool             { return false }
func (o Object) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: OBJECT} }
func (o Object) IsString() bool            { return false }
func (o Object) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: OBJECT} }
func (o Object) IsNumber() bool            { return false }
func (o Object) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: OBJECT} }
func (o Object) IsBool() bool              { return false }
func (o Object) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: OBJECT} }
func (o Object) IsNull() bool              { return false }
func (o Object) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: OBJECT} }

func (a Array) IsObject() bool            { return false }
func (a Array) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: ARRAY} }
func (a Array) IsArray() bool             { return true }
func (a Array) AsArray() (Array, error)   { return a, nil }
func (a Array) IsString() bool            { return false }
func (a Array) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: ARRAY} }
func (a Array) IsNumber() bool            { return false }
func (a Array) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: ARRAY} }
func (a Array) IsBool() bool              { return false }
func (a Array) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: ARRAY} }
func (a Array) IsNull() bool              { return false }
func (a Array) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: ARRAY} }

func (s String) IsObject() bool            { return false }
func (s String) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: STRING} }
func (s String) IsArray() bool             { return false }
func (s String) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: STRING} }
func (s String) IsString() bool            { return true }
func (s String) AsString() (String, error) { return s, nil }
func (s String) IsNumber() bool            { return false }
func (s String) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: STRING} }
func (s String) IsBool() bool              { return false }
func (s String) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: STRING} }
func (s String) IsNull() bool              { return false }
func (s String) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: STRING} }

func (n Number) IsObject() bool            { return false }
func (n Number) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: NUMBER} }
func (n Number) IsArray() bool             { return false }
func (n Number) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: NUMBER} }
func (n Number) IsString() bool            { return false }
func (n Number) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: NUMBER} }
func (n Number) IsNumber() bool            { return true }
func (n Number) AsNumber() (Number, error) { return n, nil }
func (n Number) IsBool() bool              { return false }
func (n Number) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: NUMBER} }
func (n Number) IsNull() bool              { return false }
func (n Number) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: NUMBER} }

func (b Bool) IsObject() bool            { return false }
func (b Bool) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: BOOL} }
func (b Bool) IsArray() bool             { return false }
func (b Bool) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: BOOL} }
func (b Bool) IsString() bool            { return false }
func (b Bool) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: BOOL} }
func (b Bool) IsNumber() bool            { return false }
func (b Bool) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: BOOL} }
func (b Bool) IsBool() bool              { return true }
func (b Bool) AsBool() (Bool, error)     { return b, nil }
func (b Bool) IsNull() bool              { return false }
func (b Bool) AsNull() (Null, error)     { return nil, ErrAsValue{Not: NULL, But: BOOL} }

func (n Null) IsObject() bool            { return false }
func (n Null) AsObject() (Object, error) { return nil, ErrAsValue{Not: OBJECT, But: NULL} }
func (n Null) IsArray() bool             { return false }
func (n Null) AsArray() (Array, error)   { return nil, ErrAsValue{Not: ARRAY, But: NULL} }
func (n Null) IsString() bool            { return false }
func (n Null) AsString() (String, error) { return "", ErrAsValue{Not: STRING, But: NULL} }
func (n Null) IsNumber() bool            { return false }
func (n Null) AsNumber() (Number, error) { return 0, ErrAsValue{Not: NUMBER, But: NULL} }
func (n Null) IsBool() bool              { return false }
func (n Null) AsBool() (Bool, error)     { return false, ErrAsValue{Not: BOOL, But: NULL} }
func (n Null) IsNull() bool              { return true }
func (n Null) AsNull() (Null, error)     { return n, nil }
