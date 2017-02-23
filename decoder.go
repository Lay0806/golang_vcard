// Copyright (C) 2012 Laurent Le Goff
// Copyright (C) 2015 Constantin Schomburg <me@cschomburg.com>

package golib_vcard

import (
	"bytes"
	"errors"
	"io"
	"text/scanner"
)

// A Decoder reads Directory Information Blocks from an input stream.
type Decoder struct {
	scan        *scanner.Scanner
	nextProfile string
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	var s scanner.Scanner
	s.Init(r)
	return &Decoder{&s, ""}
}

// ReadContentLine reads the next content line and returns it.
func (dec *Decoder) ReadContentLine() (*ContentLine, error) {
	dec.skipWhitespace()
	if dec.scan.Peek() == scanner.EOF {
		return nil, io.EOF
	}
	group, name := dec.readGroupName()
	params := make(map[string]Value)
	if dec.scan.Peek() == ';' {
		dec.scan.Next()
		params = dec.readParameters()
	}
	dec.scan.Next()
	value := dec.readValues(name)
	return &ContentLine{group, name, params, value}, nil
}

func (dec *Decoder) expectSingleValue(name string) (string, error) {
	cl, err := dec.ReadContentLine()
	if err != nil {
		return "", err
	}
	if cl.Name != name {
		return "", errors.New("expected " + name + ", not " + cl.Name)
	}
	return cl.Value[0][0], nil
}

// ReadObject reads the next object block and returns it.
func (dec *Decoder) ReadObject() (o *Object, err error) {
	o = &Object{}
	if dec.nextProfile == "" {
		o.Profile, err = dec.expectSingleValue("BEGIN")
		if err != nil {
			return o, err
		}
	} else {
		o.Profile = dec.nextProfile
		dec.nextProfile = ""
	}

	for {
		cl, err := dec.ReadContentLine()
		if err != nil {
			return o, err
		}
		if cl.Name == "BEGIN" {
			dec.nextProfile = cl.Value[0][0]
			comp, err := dec.ReadObject()
			if err != nil {
				return o, err
			}
			o.Objects = append(o.Objects, comp)
		}
		if cl.Name == "END" {
			if cl.Value[0][0] != o.Profile {
				return o, errors.New("unexpected END:" + cl.Value[0][0] + ", expected END:" + o.Profile)
			}
			break
		}
		if err != nil {
			return o, err
		}
		o.Properties = append(o.Properties, cl)
	}
	return o, nil
}

func (dec *Decoder) skipWhitespace() {
	c := dec.scan.Peek()
	for c == ' ' || c == '\n' || c == '\r' || c == '\t' {
		dec.scan.Next()
		c = dec.scan.Peek()
	}
}

func (dec *Decoder) readGroupName() (group, name string) {
	c := dec.scan.Peek()
	var buf []rune
	for c != scanner.EOF {
		if c == '.' {
			group = string(buf)
			buf = []rune{}
		} else if c == ';' || c == ':' {
			name = string(buf)
			return
		} else {
			buf = append(buf, c)
		}
		dec.scan.Next()
		c = dec.scan.Peek()
	}
	return
}

func (dec *Decoder) readValue(stopOnEquals bool) string {
	c := dec.scan.Peek()
	var buf []rune
	escape := false
	for c != scanner.EOF {
		if c == '\n' {
			la := dec.scan.Peek()
			if la != ' ' && la != '\t' {
				return string(buf)
			} else {
				// unfold
				for c == ' ' || c == '\t' {
					dec.scan.Next()
					c = dec.scan.Peek()
				}
			}
		}
		if c == '\\' {
			escape = true
			dec.scan.Next()
		} else if escape {
			if c == 'n' || c == 'N' {
				c = '\n'
			}
			buf = append(buf, c)
			escape = false
			dec.scan.Next()
		} else if c == ',' || c == ';' || c == ':' {
			return string(buf)
		} else if stopOnEquals && c == '=' {
			return string(buf)
		} else if c != '\n' && c != '\r' {
			buf = append(buf, c)
			dec.scan.Next()
		}
		c = dec.scan.Peek()
	}
	return string(buf)
}

func (dec *Decoder) readParameters() (params map[string]Value) {
	c := dec.scan.Peek()
	var buf []rune
	var name string
	var value string
	quoted := false
	params = make(map[string]Value)
	var values Value

	for c != scanner.EOF {
		if c == ',' && !quoted {
			values = append(values, string(buf))
			buf = []rune{}
		} else if c == ';' || c == ':' {
			if name == "" {
				name = string(buf)
			} else {
				value = string(buf)
			}
			if name != "" {
				values = append(values, value)
				if _, ok := params[name]; ok {
					params[name] = append(params[name], values...)
				} else {
					params[name] = values
				}
			}
			if c == ':' {
				return
			}
			buf = []rune{}
			values = Value{}
			name = ""
			value = ""
		} else if c == '=' {
			name = string(buf)
			buf = []rune{}
		} else if c == '"' {
			quoted = !quoted
		} else {
			buf = append(buf, c)
		}
		dec.scan.Next()
		c = dec.scan.Peek()
	}
	return
}

func (dec *Decoder) readValues(name string) (value StructuredValue) {

	if name == "END" {
		var valu Value
		var buff []rune
		cc := dec.scan.Next()
		for cc != scanner.EOF && cc != '\n' && cc != '\r' {
			buff = append(buff, cc)
			cc = dec.scan.Next()
		}
		valu = append(valu, string(buff))
		value = append(value, valu)
	} else {
		c := dec.scan.Next()
		var buf []rune
		escape := false
		var val Value
		for c != scanner.EOF {
			if c == '\n' {
				la := dec.scan.Peek()
				if la != 32 && la != 9 {
					// return
					if len(buf) > 0 {
						val = append(val, string(buf))
					}
					value = append(value, val)
					return
				} else {
					// unfold
					c = dec.scan.Next()
					for c == 32 || c == 9 {
						c = dec.scan.Next()
					}
				}
			}
			if c == '\\' {
				escape = true
			} else if escape {
				if c == 'n' || c == 'N' {
					c = '\n'
				}
				buf = append(buf, c)
				escape = false
			} else if c == ',' {
				if len(buf) > 0 {
					val = append(val, string(buf))
					buf = []rune{}
				}
			} else if c == ';' {
				if len(buf) > 0 {
					val = append(val, string(buf))
					buf = []rune{}
				}
				value = append(value, val)
				val = Value{}
			} else if c != '\n' && c != '\r' {
				buf = append(buf, c)
			}
			c = dec.scan.Next()
		}
	}

	return
}

// Decode reads the next object and stores it in the value pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion into
// a Go value.
func (dec *Decoder) Decode(v interface{}) error {
	o, err := dec.ReadObject()
	if err != nil {
		return err
	}
	return FromObject(v, o)
}

// Unmarshal parses a single directory block from data and stores the result in
// the value pointed to by v.
//
// To unmarshal an object into a struct, Unmarshal matches incoming object
// properties and components to the keys used by Marshal (either the struct
// field name or its tag), accepting a case-insensitive match.
//
// A property content line can be unmarshalled into a string for a single
// value, a string slice for a value list (delimited by commas) or a struct for
// structured values. Umarshalling into a struct first maps all struct fields
// with tag ",param" to the respective parameter values and fills the remaining
// fields in their index order with the respective semicolon-delimited value
// components.
func Unmarshal(data []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewReader(data))
	return dec.Decode(v)
}
