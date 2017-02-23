// Package vdir implements RFC 2425 directory encoding with support for vCard
// and iCalendar profiles.
//
// The directory blocks can be mapped to an arbitrary Go Value which is
// described in the Marshal and Unmarshal functions. In addition, the package
// provides conversions to/from a generalized Object representation which
// allows for detailed access and easier modification.
package golib_vcard

// Object is a generic Directory Information Block.
type Object struct {
	Profile    string
	Properties []*ContentLine
	Objects    []*Object
}

// PropertyMap returns all properties indexed by their name.
func (o *Object) PropertyMap() map[string][]*ContentLine {
	props := make(map[string][]*ContentLine)
	for _, cl := range o.Properties {
		props[cl.Name] = append(props[cl.Name], cl)
	}
	return props
}

// ContentLine is a single named line with named parameters and values.
type ContentLine struct {
	Group, Name string
	Params      map[string]Value
	Value       StructuredValue
}

type StructuredValue []Value
type Value []string

// GetText returns the first value of the first component.
func (v StructuredValue) GetText() string {
	if len(v) > 0 && len(v[0]) > 0 {
		return v[0][0]
	}
	return ""
}

// GetText returns the first value.
func (v Value) GetText() string {
	if len(v) > 0 {
		return v[0]
	}
	return ""
}
