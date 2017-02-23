package golib_vcard

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// An encoder writes Directory Information Blocks to an input stream.
type Encoder struct {
	writer io.Writer
	err    error
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w, nil}
}

// WriteObject writes an object to the stream.
func (enc *Encoder) WriteObject(o *Object) error {
	if o.Profile == "" {
		return errors.New("No profile set.")
	}
	enc.WriteContentLine(&ContentLine{"", "BEGIN", nil, StructuredValue{Value{o.Profile}}})
	for _, p := range o.Properties {
		enc.WriteContentLine(p)
	}
	for _, so := range o.Objects {
		enc.WriteObject(so)
	}
	enc.WriteContentLine(&ContentLine{"", "END", nil, StructuredValue{Value{o.Profile}}})
	return enc.err
}

func (enc *Encoder) writeString(s string) error {
	if enc.err != nil {
		return enc.err
	}
	_, enc.err = io.WriteString(enc.writer, s)
	return enc.err
}

// WriteContentLine writes a single content line to the stream.
func (enc *Encoder) WriteContentLine(cl *ContentLine) error {
	if cl.Group != "" {
		enc.writeString(cl.Group)
		enc.writeString(".")
	}
	enc.writeString(cl.Name)

	if cl.Params != nil {
		for key, values := range cl.Params {
			enc.writeString(";")
			enc.writeString(key)
			if len(values) > 0 {
				enc.writeString("=")
				for vi := 0; vi < len(values); vi++ {
					enc.writeParamValue(values[vi])
					if vi+1 < len(values) {
						enc.writeString(",")
					}
				}
			}
		}
	}
	enc.writeString(":")
	for si := 0; si < len(cl.Value); si++ {
		for vi := 0; vi < len(cl.Value[si]); vi++ {
			enc.writeValue(cl.Value[si][vi])
			if vi+1 < len(cl.Value[si]) {
				enc.writeString(",")
			}
		}
		if si+1 < len(cl.Value) {
			enc.writeString(";")
		}
	}
	enc.writeString("\r\n")
	return enc.err
}

func (enc *Encoder) writeValue(v string) error {
	i := 0
	for _, c := range v {
		if i == 76 {
			enc.writeString("\n  ")
			i = 0
		}
		var e string
		switch c {
		case '\r':
			e = `\r`
		case '\n':
			e = `\n`
		case ';':
			e = `\;`
		case ',':
			e = `\,`
		default:
			e = string(c)
		}
		enc.writeString(e)
		i++
	}
	return enc.err
}

func (enc *Encoder) writeParamValue(v string) error {
	quoted := strings.ContainsAny(v, `:;,`)
	if quoted {
		v = `"` + v + `"`
	}
	v = strings.Replace(v, "\n", `\n`, -1)
	return enc.writeString(v)
}

// Encode writes the encoded block of v to the stream.
//
// See the documentation for Marshal for details about the conversion
// of Go values.
func (enc *Encoder) Encode(v interface{}) error {
	o := &Object{}
	if err := ToObject(v, o); err != nil {
		return err
	}
	return enc.WriteObject(o)
}

// Marshal returns the encoding of v.
//
// Marshal expects a struct as value v and traverses it recursively,
// mapping its fields to properties and component objects.
//
// Each struct field becomes a property of the object unless the fields tag is
// "-" or the field is empty. The tag value is the property name, otherwise
// the uppercase struct field name is used.
//
// A special field tagged with ",profile" defines the profile type of the block,
// for example a tag "vcard,profile" designates the struct to be of type VCARD.
//
// Fields of type string are mapped to single values, string slices to a
// comma-delimited value list.
//
// Fields that contain a struct are stored as a structured property with
// optional parameters and components. If a struct fields tag contains a "param"
// option as second value, it is stored as a parameter of the property. Untagged
// fields are stored as semicolon-delimited component values based on their order
// of appearance in the struct.
//
// Fields that have a tag with option "objects" as second value are converted
// to a new inner BEGIN:PROFILE-END object block.
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	err := enc.Encode(v)
	return b.Bytes(), err
}
