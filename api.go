package xobj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/worldiety/jsonml"
	"strconv"
)

// StrList is used instead of a slice, to be directly compatible with gomobile.
type StrList interface {
	// Size returns the amount of entries in the list
	Size() int

	// Get returns the string at the given index. Panics if out of bounds.
	Get(idx int) string
}

// An Obj allows a key/value oriented access with various comfort functions.
type Obj interface {
	// Keys returns a list of declared string names
	Keys() StrList

	// Get returns the generic value associated with the key. You cannot distinguish a null value from the non-existence
	// of the key. Use #Has() method to validate. The method is discarded when used with gomobile.
	Get(name string) interface{}

	// Put inserts just a generic value and associates it with the key. The method is discarded when used with gomobile.
	Put(name string, value interface{}) Obj

	// Remove deletes the key/value combination from this object. Returns the object for a builder pattern.
	// Returning the Obj avoids discarding the method in gomobile.
	Remove(name string) Obj

	// Has returns true, if the given key is available in this object.
	Has(name string) bool

	// AsInt64 tries to convert the associated value into an int64, otherwise returns an error
	AsInt64(name string) (int64, error)

	// OptInt64 tries to convert the associated value into an int64 or returns the fallback
	OptInt64(name string, fallback int64) int64

	// PutInt64 removes the existing field and replaces its value. Returns the object for a builder pattern.
	PutInt64(name string, value int64) Obj

	// AsBool tries to convert the associated value into a boolean, otherwise returns an error
	AsBool(name string) (bool, error)

	// OptBool tries to convert the associated value into a bool or returns the fallback
	OptBool(name string, fallback bool) bool

	// PutBool removes the existing field and replaces its value. Returns the object for a builder pattern.
	PutBool(name string, value bool) Obj

	// AsFloat64 tries to convert the associated value into an float64, otherwise returns an error
	AsFloat64(name string) (float64, error)

	// OptFloat64 tries to convert the associated value into a float64 or returns the fallback
	OptFloat64(name string, fallback float64) float64

	// PutFloat64 removes the existing field and replaces its value. Returns the object for a builder pattern.
	PutFloat64(name string, value float64) Obj

	// AsString tries to convert the associated value into an string, otherwise returns an error
	AsString(name string) (string, error)

	// OptString tries to convert the associated value into a string or returns the fallback
	OptString(name string, fallback string) string

	// PutString removes the existing field and replaces its value. Returns the object for a builder pattern.
	PutString(name string, value string) Obj

	// AsObject returns the value as an Obj, if the type matches, otherwise returns an error
	AsObject(name string) (Obj, error)

	// OptObject returns the value as an Obj, if the type matches, otherwise creates(!) a new
	// Object, puts it to the named field and returns the new instance.
	OptObject(name string) Obj

	// PutObject removes the existing field and replaces its value. Returns the object for a builder pattern.
	PutObject(name string, value Obj) Obj

	// AsObject returns the value as an Arr, if the type matches, otherwise returns an error
	AsArray(name string) (Arr, error)

	// OptArray returns the value as an Array, if the type matches, otherwise creates(!) a new
	// Array, puts it to the named field and returns the new instance.
	OptArray(name string) Arr

	// PutArray removes the existing field and replaces its value. Returns the object for a builder pattern.
	PutArray(name string, value Arr) Obj

	// String provides the stringer interface, which returns a compact JSON serialization
	String() string
}

// An Arr is typed accessor for index-based and ordered access of data.
type Arr interface {
	// Size returns the amount of entries in this Array
	Size() int

	// Get returns the generic value associated with the index.
	// The method is discarded when used with gomobile.
	Get(idx int) interface{}

	// Put replaces the value at the given index.
	// The method is discarded when used with gomobile.
	Put(idx int, value interface{}) Arr

	// Remove deletes the value at the given index. Returns the Arr for a builder pattern.
	// Returning the Obj avoids discarding the method in gomobile.
	Remove(idx int) Arr

	// AsInt64 tries to convert the associated value into an int64, otherwise returns an error
	AsInt64(idx int) (int64, error)

	// OptInt64 tries to convert the associated value into an int64 or returns the fallback
	OptInt64(idx int, fallback int64) int64

	// PutInt64 replaces the value at the given index. Returns the array for a builder pattern. Panics if idx is out
	// of bounds.
	PutInt64(idx int, value int64) Arr

	// AddInt64 appends the value and returns the array.
	AddInt64(value int64) Arr

	// AsBool tries to convert the associated value into a boolean, otherwise returns an error
	AsBool(idx int) (bool, error)

	// OptBool tries to convert the associated value into a bool or returns the fallback
	OptBool(idx int, fallback bool) bool

	// PutBool replaces the value at the given index. Returns the array for a builder pattern. Panics if idx is out
	// of bounds.
	PutBool(idx int, value bool) Arr

	// AddBool appends the value and returns the array.
	AddBool(value bool) Arr

	// AsFloat64 tries to convert the associated value into an float64, otherwise returns an error
	AsFloat64(idx int) (float64, error)

	// OptFloat64 tries to convert the associated value into a float64 or returns the fallback
	OptFloat64(idx int, fallback float64) float64

	// PutFloat64 replaces the value at the given index. Returns the array for a builder pattern. Panics if idx is out
	// of bounds.
	PutFloat64(idx int, value float64) Arr

	// AddFloat64 appends the value and returns the array.
	AddFloat64(value float64) Arr

	// AsString tries to convert the associated value into an string, otherwise returns an error
	AsString(idx int) (string, error)

	// OptString tries to convert the associated value into a string or returns the fallback
	OptString(idx int, fallback string) string

	// PutString replaces the value at the given index. Returns the array for a builder pattern. Panics if idx is out
	// of bounds.
	PutString(idx int, value string) Arr

	// AddString appends the value and returns the array.
	AddString(value string) Arr

	// AsObject returns the value as an Obj, if the type matches, otherwise returns an error
	AsObject(idx int) (Obj, error)

	// OptObject returns the value as an Object, if the type matches, otherwise creates(!) a new
	// Object, puts it to the index and returns the new instance. Panics if index is out of bounds.
	OptObject(idx int) Obj

	// PutObject replaces the value at the given index. Returns the array for a builder pattern. Panics if idx is out
	// of bounds.
	PutObject(idx int, value Obj) Arr

	// AddInt64 appends the value and returns the array.
	AddObject(value Obj) Arr

	// AsObject returns the value as an Arr, if the type matches, otherwise returns an error
	AsArray(idx int) (Arr, error)

	// OptArray returns the value as an Array, if the type matches, otherwise creates(!) a new
	// Array, puts it to the index and returns the new instance. Panics if index is out of bounds.
	OptArray(idx int) Arr

	// PutArray replaces the value at the given index. Returns the array for a builder pattern. Panics if idx is out
	// of bounds.
	PutArray(idx int, value Arr) Arr

	// AddArray appends the value and returns the array.
	AddArray(value Arr) Arr

	// String provides the stringer interface, which returns a compact JSON serialization
	String() string
}

// ToString converts anything to a string
func ToString(any interface{}) string {

	switch t := any.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64)
	case int64:
		return strconv.FormatInt(t, 10)
	case bool:
		return strconv.FormatBool(t)
	}
	return fmt.Sprintf("%v", any)
}

// NewObj creates a new instance of Object
func NewObj() Obj {
	return Object{}
}

// NewArr creates a new instance of Array
func NewArr() Arr {
	return &Array{}
}

// Parse tries to parse the given bytes either as XML or as JSON.
// If data looks like XML, a jsonml transformation is applied, which
// is available in the field 'xml'.
// If data represents an array, it is wrapped automatically into
// an object, using the field name "array".
func Parse(data []byte) (Obj, error) {
	obj := Object{}
	err := json.Unmarshal(data, &obj)
	if err == nil {
		return obj, nil
	}

	// failed with json object, try to parse as array
	arr := &Array{}
	err = json.Unmarshal(data, &arr)
	if err == nil {
		obj.Put("array", arr)
		return obj, nil
	}

	// failed with json array, try to parse as xml
	slice, err := jsonml.ToJSON(true, bytes.NewReader(data))
	if err == nil {
		obj.Put("xml", slice)
		return obj, err
	}

	// failed, give up
	return obj, err
}
