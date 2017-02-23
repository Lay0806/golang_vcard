package golib_vcard

import (
	"errors"
	"reflect"
	"strings"
)

var (
	ErrEmpty = errors.New("Empty")
)

// ToObject converts a struct v into an intermediate Object that
// can be written by an encoder.
//
// See the documentation for Marshal for details about the conversion of a
// Go Value.
func ToObject(v interface{}, o *Object) error {
	if v == nil {
		return errors.New("Can not marshal nil value.")
	}
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	switch rv.Kind() {
	case reflect.Struct:
		typ := rv.Type()
		for i := 0; i < typ.NumField(); i++ {
			rvi := rv.Field(i)
			name, opt := fieldToProp(typ.Field(i))
			if name == "-" {
				continue
			}
			switch opt {
			case "profile":
				if v := rvi.String(); v != "" {
					o.Profile = v
				} else if name != "" {
					o.Profile = name
				}
			case "object":
				if rvi.Kind() == reflect.Slice {
					if rvi.IsNil() {
						continue
					}
					for j := 0; j < rvi.Len(); j++ {
						comp := &Object{}
						if err := ToObject(rvi.Index(j), comp); err != nil {
							return err
						}
						o.Objects = append(o.Objects, comp)
					}
				} else {
					comp := &Object{}
					if err := ToObject(rvi, comp); err != nil {
						return err
					}
					o.Objects = append(o.Objects, comp)
				}
			default:
				if rvi.Kind() == reflect.Slice && rvi.Type().Elem().Kind() != reflect.String {
					for j := 0; j < rvi.Len(); j++ {
						cl := &ContentLine{Name: name}
						if err := toContentLine(rvi.Index(j), cl); err != nil {
							if err == ErrEmpty {
								continue
							}
							return err
						}
						o.Properties = append(o.Properties, cl)
					}
				} else {
					cl := &ContentLine{Name: name}
					if err := toContentLine(rv.Field(i), cl); err != nil {
						if err == ErrEmpty {
							continue
						}
						return err
					}
					o.Properties = append(o.Properties, cl)
				}
			}
		}
	default:
		return errors.New("Cannot marshal " + rv.Type().String() + " into object")
	}

	return nil
}

func toContentLine(rv reflect.Value, cl *ContentLine) error {
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	switch rv.Kind() {
	case reflect.Struct:
		typ := rv.Type()
		for i := 0; i < typ.NumField(); i++ {
			name, opt := fieldToProp(typ.Field(i))
			if name == "-" {
				continue
			}
			v, err := toValue(rv.Field(i))
			if err != nil {
				return err
			}
			if opt == "param" {
				if len(v) == 0 || v[0] == "" {
					continue
				}
				if cl.Params == nil {
					cl.Params = make(map[string]Value)
				}
				cl.Params[name] = v
			} else {
				cl.Value = append(cl.Value, v)
			}
		}
	default:
		v, err := toValue(rv)
		if err != nil {
			return err
		}
		if len(v) == 0 || v[0] == "" {
			return ErrEmpty
		}
		cl.Value = StructuredValue{v}
	}
	return nil
}

func toValue(rv reflect.Value) (Value, error) {
	switch rv.Kind() {
	case reflect.String:
		return Value{rv.String()}, nil
	case reflect.Slice:
		if rv.Type().Elem().Kind() != reflect.String {
			return nil, errors.New("Cannot marshal " + rv.Type().String() + " into value")
		}
		v := Value{}
		for i := 0; i < rv.Len(); i++ {
			v = append(v, rv.Index(i).String())
		}
		return v, nil
	}
	return nil, errors.New("Cannot marshal " + rv.Type().String() + " into value")
}

func fieldToProp(f reflect.StructField) (name, options string) {
	tag := strings.Split(f.Tag.Get("vdir"), ",")
	name = tag[0]
	if len(tag) > 1 {
		options = tag[1]
	}
	if name == "" {
		name = f.Name
	}
	return strings.ToUpper(name), options
}

// FromObject converts an intermediate Object into a struct.
//
// See the documentation for Unmarshal for details about the conversion of into
// a Go Value.
func FromObject(v interface{}, o *Object) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr {
		return errors.New("Cannot unmarshal oject into non-pointer " + rv.Type().String())
	}
	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}

	props := o.PropertyMap()
	switch rv.Elem().Kind() {
	case reflect.Struct:
		typ := rv.Type()
		for i := 0; i < typ.Elem().NumField(); i++ {
			rvi := rv.Elem().Field(i)
			name, opt := fieldToProp(typ.Elem().Field(i))
			if name == "-" {
				continue
			}
			switch opt {
			case "profile":
				if rvi.Kind() != reflect.String {
					return errors.New("Cannot unmarshal profile into " + rvi.Type().String())
				}
				rvi.SetString(o.Profile)
			case "object":
				if rvi.Kind() != reflect.Slice {
					return errors.New("Cannot unmarshal object into " + rv.Type().String())
				}
				for _, so := range o.Objects {
					if so.Profile != name {
						continue
					}
					rvii := reflect.New(rvi.Type().Elem())
					if err := FromObject(rvii.Interface(), so); err != nil {
						return err
					}
					rvi.Set(reflect.Append(rvi, reflect.Indirect(rvii)))
				}
			default:
				cls := props[name]
				if cls == nil || len(cls) == 0 {
					continue
				}
				if rvi.Kind() == reflect.Slice && rvi.Type().Elem().Kind() != reflect.String {
					for _, cl := range cls {
						rvii := reflect.New(rvi.Type().Elem())
						if err := fromContentLine(rvii, cl); err != nil {
							return err
						}
						rvi.Set(reflect.Append(rvi, reflect.Indirect(rvii)))
					}
				} else {
					if err := fromContentLine(rvi.Addr(), cls[0]); err != nil {
						return err
					}
				}
			}
		}
	default:
		return errors.New("Cannot unmarshal object into " + rv.Type().String())
	}
	return nil
}

func fromContentLine(rv reflect.Value, cl *ContentLine) error {
	if rv.Kind() != reflect.Ptr {
		return errors.New("Cannot unmarshal property into non-pointer " + rv.Type().String())
	}
	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}

	switch rv.Elem().Kind() {
	case reflect.Struct:
		typ := rv.Type()
		vi := 0
		for i := 0; i < typ.Elem().NumField(); i++ {
			name, opt := fieldToProp(typ.Elem().Field(i))
			if name == "-" {
				continue
			}
			if opt == "param" {
				if v, ok := cl.Params[name]; ok {
					if err := fromValue(rv.Elem().Field(i).Addr(), v); err != nil {
						return err
					}
				}
			} else {
				if len(cl.Value) > vi {
					if err := fromValue(rv.Elem().Field(i).Addr(), cl.Value[vi]); err != nil {
						return err
					}
					vi++
				}
			}
		}
		return nil
	default:
		if len(cl.Value) > 0 {
			if err := fromValue(rv, cl.Value[0]); err != nil {
				return err
			}
		}
	}
	return nil
}

func fromValue(rv reflect.Value, v Value) error {
	if rv.Kind() != reflect.Ptr {
		return errors.New("Cannot unmarshal value into non-pointer " + rv.Type().String())
	}
	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}

	switch rv.Elem().Kind() {
	case reflect.String:
		rv.Elem().SetString(v.GetText())
		return nil
	case reflect.Slice:
		if rv.Type().Elem().Elem().Kind() != reflect.String {
			break
		}
		for _, vv := range v {
			rv.Elem().Set(reflect.Append(rv.Elem(), reflect.ValueOf(vv)))
		}
		return nil
	}
	return errors.New("Cannot unmarshal value into " + rv.Type().String())
}
