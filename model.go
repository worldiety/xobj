package xobj

import (
	"encoding/json"
	"fmt"
	"math"
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

	// IsNull returns false if either the Obj has no such field or its field value is null/nil. Use #Has() to
	// differentiate between the cases.
	IsNull(name string) bool

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

	// IsNull checks if the value at the given index is null. It also returns false, if idx is out of bounds.
	IsNull(idx int) bool

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
	if any == nil {
		return ""
	}
	switch t := any.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64)
	case int64:
		return strconv.FormatInt(t, 10)
	case bool:
		return strconv.FormatBool(t)
	case fmt.Stringer:
		return t.String()
	case []interface{}:
		data, err := json.Marshal(t)
		if err != nil {
			logger.Info(Fields{"msg": "failed to marshal", "err": err})
			return fmt.Sprintf("%v", any)
		}
		return string(data)
	case map[string]interface{}:
		data, err := json.Marshal(t)
		if err != nil {
			logger.Info(Fields{"msg": "failed to marshal", "err": err})
			return fmt.Sprintf("%v", any)
		}
		return string(data)

	case *[]interface{}:
		return ToString(*t)
	}
	return fmt.Sprintf("%v", any)
}

// ToFloat64 converts anything to a float
func ToFloat64(any interface{}) float64 {
	if any == nil {
		return math.NaN()
	}
	switch t := any.(type) {

	case float64:
		return t
	case int64:
		return float64(t)
	case uint64:
		return float64(t)
	case int32:
		return float64(t)
	case uint32:
		return float64(t)
	case int16:
		return float64(t)
	case uint16:
		return float64(t)
	case uint8:
		return float64(t)
	case int8:
		return float64(t)
	case bool:
		if t {
			return 1
		} else {
			return 0
		}
	case string:
		return parseStr(t)
	case fmt.Stringer:
		return parseStr(t.String())
	default:
		return parseStr(fmt.Sprintf("%v", any))

	}

}

func parseStr(str string) float64 {
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return f
	}
	return math.NaN()
}

// NewObj creates a new instance of Object
func NewObj() Obj {
	return Object{}
}

// NewArr creates a new instance of Array
func NewArr() Arr {
	return &Array{}
}

// UnwrapObj allocates new maps and slices for Obj and Arr instances. We do not cast recursively
// because it would change internal types, values and pointers while running causing all sorts of bugs.
func UnwrapObj(obj Obj) map[string]interface{} {
	res := make(map[string]interface{})
	keys := obj.Keys()
	for i := 0; i < keys.Size(); i++ {
		k := keys.Get(i)
		c, err := obj.AsObject(k)
		if err == nil {
			res[keys.Get(i)] = UnwrapObj(c)
			continue
		}

		a, err := obj.AsArray(k)
		if err == nil {
			res[keys.Get(i)] = UnwrapArr(a)
			continue
		}
		res[k] = obj.Get(k)
	}
	return res
}

// UnwrapArr allocates new maps and slices for Obj and Arr instances. We do not cast recursively
// because it would change internal types, values and pointers while running causing all sorts of bugs.
func UnwrapArr(arr Arr) []interface{} {
	res := make([]interface{}, arr.Size())
	for i := 0; i < arr.Size(); i++ {
		c, err := arr.AsObject(i)
		if err == nil {
			res[i] = UnwrapObj(c)
			continue
		}

		a, err := arr.AsArray(i)
		if err == nil {
			res[i] = UnwrapArr(a)
			continue
		}
		res[i] = arr.Get(i)
	}
	return res
}
