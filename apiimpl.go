package xobj

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

//=

var _ StrList = (StringList)(nil)

// A StringList is just a simple slice of strings, so you can perform an efficient type conversion.
type StringList []string

func (s StringList) Size() int {
	return len(s)
}

func (s StringList) Get(idx int) string {
	return s[idx]
}

//=

var _ Obj = (Object)(nil)

// An Object is just a simple map of interfaces, so you can perform an efficient type conversion.
type Object map[string]interface{}

func (o Object) Keys() StrList {
	res := StringList{}
	for k := range o {
		res = append(res, k)
	}
	return res
}

func (o Object) Get(name string) interface{} {
	return o[name]
}

func (o Object) Put(name string, value interface{}) Obj {
	o[name] = value
	return o
}

func (o Object) Remove(name string) Obj {
	delete(o, name)
	return o
}

func (o Object) Has(name string) bool {
	_, has := o[name]
	return has
}

// AsInt64 uses a lot of type assertions to optimize performance
func (o Object) AsInt64(name string) (int64, error) {
	v, ok := o[name]
	if !ok {
		return 0, unknownFieldName(name)
	}
	return asInt64(v)
}

func (o Object) OptInt64(name string, fallback int64) int64 {
	v, err := o.AsInt64(name)
	if err != nil {
		return fallback
	}
	return v
}

func (o Object) PutInt64(name string, value int64) Obj {
	o[name] = value
	return o
}

func (o Object) AsBool(name string) (bool, error) {
	v, ok := o[name]
	if !ok {
		return false, unknownFieldName(name)
	}
	return asBool(v)
}

func (o Object) OptBool(name string, fallback bool) bool {
	v, err := o.AsBool(name)
	if err != nil {
		return fallback
	}
	return v
}

func (o Object) PutBool(name string, value bool) Obj {
	o[name] = value
	return o
}

// AsFloat64 uses a lot of type assertions to optimize performance
func (o Object) AsFloat64(name string) (float64, error) {
	v, ok := o[name]
	if !ok {
		return 0, unknownFieldName(name)
	}
	return asFloat64(v)
}

func (o Object) OptFloat64(name string, fallback float64) float64 {
	v, err := o.AsFloat64(name)
	if err != nil {
		return fallback
	}
	return v
}

func (o Object) PutFloat64(name string, value float64) Obj {
	o[name] = value
	return o
}

func (o Object) AsString(name string) (string, error) {
	v, ok := o[name]
	if !ok {
		return "", unknownFieldName(name)
	}
	return asString(v)
}

func (o Object) OptString(name string, fallback string) string {
	v, err := o.AsString(name)
	if err != nil {
		return fallback
	}
	return v
}

func (o Object) PutString(name string, value string) Obj {
	o[name] = value
	return o
}

func (o Object) AsObject(name string) (Obj, error) {
	v, ok := o[name]
	if !ok {
		return nil, unknownFieldName(name)
	}
	if obj, ok := v.(Obj); ok {
		return obj, nil
	}
	return nil, fmt.Errorf("%s is not an object (%v)", name, reflect.TypeOf(v))
}

func (o Object) OptObject(name string) Obj {
	v, err := o.AsObject(name)
	if err != nil {
		v = Object{}
		o.PutObject(name, v)
		return v
	}
	return v
}

func (o Object) PutObject(name string, value Obj) Obj {
	o[name] = value
	return o
}

func (o Object) AsArray(name string) (Arr, error) {
	v, ok := o[name]
	if !ok {
		return nil, unknownFieldName(name)
	}
	if arr, ok := v.(Arr); ok {
		return arr, nil
	}
	if arr, ok := v.(*[]interface{}); ok {
		return (*Array)(arr), nil
	}
	//perform a replacement to a slice pointer, so that the thing can exchange the slice struct, e.g. for appending
	if arr, ok := v.([]interface{}); ok {
		tmp := Array(arr)
		o[name] = &tmp
		return &tmp, nil
	}
	return nil, fmt.Errorf("value in field '%s' is not an array, is type '%v'", name, reflect.TypeOf(v))
}

func (o Object) OptArray(name string) Arr {
	v, err := o.AsArray(name)
	if err != nil {
		v = &Array{}
		o.PutArray(name, v)
		return v
	}
	return v
}

func (o Object) PutArray(name string, value Arr) Obj {
	o[name] = value
	return o
}

func (o Object) String() string {
	str, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

//==

var _ Arr = (*Array)(nil)

// An Array (used with a pointer) is just a simple slice of interfaces, so you can perform an efficient type conversion.
// The pointer receiver is required, to apply changes to the underling slice reference.
type Array []interface{}

func (a *Array) Size() int {
	return len(*a)
}

func (a *Array) Get(idx int) interface{} {
	return (*a)[idx]
}

func (a *Array) Put(idx int, value interface{}) Arr {
	(*a)[idx] = value
	return a
}

func (a *Array) Remove(idx int) Arr {
	copy((*a)[idx:], (*a)[idx+1:])
	(*a)[len(*a)-1] = nil //avoid future memory leak of the last element in "hidden" slice element
	*a = (*a)[:len(*a)-1]
	return a
}

func (a *Array) AsInt64(idx int) (int64, error) {
	if idx < 0 || idx >= len(*a) {
		return 0, outOfBounds(*a, idx)
	}
	return asInt64((*a)[idx])
}

func (a *Array) OptInt64(idx int, fallback int64) int64 {
	v, err := a.AsInt64(idx)
	if err != nil {
		return fallback
	}
	return v
}

func (a *Array) PutInt64(idx int, value int64) Arr {
	(*a)[idx] = value
	return a
}

func (a *Array) AddInt64(value int64) Arr {
	*a = append(*a, value)
	return a
}

func (a *Array) AsBool(idx int) (bool, error) {
	if idx < 0 || idx >= len(*a) {
		return false, outOfBounds(*a, idx)
	}
	return asBool((*a)[idx])
}

func (a *Array) OptBool(idx int, fallback bool) bool {
	v, err := a.AsBool(idx)
	if err != nil {
		return fallback
	}
	return v
}

func (a *Array) PutBool(idx int, value bool) Arr {
	(*a)[idx] = value
	return a
}

func (a *Array) AddBool(value bool) Arr {
	*a = append(*a, value)
	return a
}

func (a *Array) AsFloat64(idx int) (float64, error) {
	if idx < 0 || idx >= len(*a) {
		return 0, outOfBounds(*a, idx)
	}
	return asFloat64((*a)[idx])
}

func (a *Array) OptFloat64(idx int, fallback float64) float64 {
	v, err := a.AsFloat64(idx)
	if err != nil {
		return fallback
	}
	return v
}

func (a *Array) PutFloat64(idx int, value float64) Arr {
	(*a)[idx] = value
	return a
}

func (a *Array) AddFloat64(value float64) Arr {
	*a = append(*a, value)
	return a
}

func (a *Array) AsString(idx int) (string, error) {
	if idx < 0 || idx >= len(*a) {
		return "", outOfBounds(*a, idx)
	}
	return asString((*a)[idx])
}

func (a *Array) OptString(idx int, fallback string) string {
	v, err := a.AsString(idx)
	if err != nil {
		return fallback
	}
	return v
}

func (a *Array) PutString(idx int, value string) Arr {
	(*a)[idx] = value
	return a
}

func (a *Array) AddString(value string) Arr {
	*a = append(*a, value)
	return a
}

func (a *Array) AsObject(idx int) (Obj, error) {
	if idx < 0 || idx >= len(*a) {
		return nil, outOfBounds(*a, idx)
	}
	if obj, ok := (*a)[idx].(Obj); ok {
		return obj, nil
	}
	if obj, ok := (*a)[idx].(map[string]interface{}); ok {
		return Object(obj), nil
	}
	return nil, fmt.Errorf("element at index %d is not an object (%v)", idx, reflect.TypeOf((*a)[idx]))
}

func (a *Array) OptObject(idx int) Obj {
	v, err := a.AsObject(idx)
	if err != nil {
		v = Object{}
		(*a)[idx] = v
		return v
	}
	return v
}

func (a *Array) PutObject(idx int, value Obj) Arr {
	(*a)[idx] = value
	return a
}

func (a *Array) AddObject(value Obj) Arr {
	*a = append(*a, value)
	return a
}

func (a *Array) AsArray(idx int) (Arr, error) {
	if idx < 0 || idx >= len(*a) {
		return nil, outOfBounds(*a, idx)
	}

	v := (*a)[idx]

	if arr, ok := v.(Arr); ok {
		return arr, nil
	}
	if arr, ok := v.(*[]interface{}); ok {
		return (*Array)(arr), nil
	}

	//perform a replacement to a slice pointer, so that the thing can exchange the slice struct, e.g. for appending
	if arr, ok := v.([]interface{}); ok {
		tmp := Array(arr)
		(*a)[idx] = &tmp
		return &tmp, nil
	}

	// hacky way of doing so, can we inspect the base type instead?
	ref := reflect.ValueOf(v)
	if ref.Kind() == reflect.Ptr {
		if ref.Elem().Kind() == reflect.Slice{
			if ref.Elem().Type().String() == "jsonml.jNode"{
				tmp := v.(*[]interface{})
				return (*Array)(tmp), nil
			}
		}
//TODO fix me by introducing an interface contract to return the correct base type, to perform the conversion
	}

	return nil, fmt.Errorf("value at index '%d' is not an array, but '%v'", idx, reflect.TypeOf(v))
}

func (a *Array) OptArray(idx int) Arr {
	v, err := a.AsArray(idx)
	if err != nil {
		fmt.Println(err)
		v = &Array{}
		(*a)[idx] = v
		return v
	}
	return v
}

func (a *Array) PutArray(idx int, value Arr) Arr {
	(*a)[idx] = value
	return a
}

func (a *Array) AddArray(value Arr) Arr {
	*a = append(*a, value)
	return a
}

func (a *Array) String() string {
	str, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

//==

func unknownFieldName(name string) error {
	return fmt.Errorf("unknown name %s", name)
}

func outOfBounds(slice []interface{}, idx int) error {
	return fmt.Errorf("out of bounds %d, having %d", idx, len(slice))
}

func asFloat64(v interface{}) (float64, error) {

	switch t := v.(type) {
	case float64:
		return t, nil
	case float32:
		return float64(t), nil
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case uint64:
		return float64(t), nil
	case bool:
		if t {
			return 1, nil
		} else {
			return 0, nil
		}
	case int8:
		return float64(t), nil
	case uint8:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case uint32:
		return float64(t), nil

	}

	return strconv.ParseFloat(ToString(v), 64)
}

func asInt64(v interface{}) (int64, error) {
	switch t := v.(type) {
	case int:
		return int64(t), nil
	case int64:
		return t, nil
	case uint64:
		return int64(t), nil
	case bool:
		if t {
			return 1, nil
		} else {
			return 0, nil
		}
	case int8:
		return int64(t), nil
	case uint8:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case uint32:
		return int64(t), nil

	}
	return strconv.ParseInt(ToString(v), 10, 64)
}

func asBool(v interface{}) (bool, error) {
	switch t := v.(type) {
	case bool:
		return t, nil
	case int:
		if t == 0 {
			return false, nil
		} else if t == 1 {
			return true, nil
		}
	}
	return strconv.ParseBool(ToString(v))
}

func asString(v interface{}) (string, error) {
	switch t := v.(type) {
	case string:
		return t, nil
	case map[string]interface{}:
		return "", fmt.Errorf("won't convert Obj to string")
	case []interface{}:
		return "", fmt.Errorf("won't convert Arr to string")
	case *[]interface{}: // we are replacing slices with pointer to slice whenever accessed
		return "", fmt.Errorf("won't convert Arr to string")
	}
	return ToString(v), nil
}
