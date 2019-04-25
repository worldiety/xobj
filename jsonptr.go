package xobj

import (
	"fmt"
	"strconv"
	"strings"
)

// A JSONPointer specifies a specific value within a JSON document.
// See https://tools.ietf.org/html/rfc6901 for the specification.
type JSONPointer = string

// A JSONPointerToken is a single element or token of a JSONPointer
type JSONPointerToken = string

// Evaluate takes the json pointer and applies it to the given json object or array.
// Returns an error if the json pointer cannot be resolved.
func Evaluate(objOrArr interface{}, ptr JSONPointer) (interface{}, error) {
	if len(ptr) == 0 {
		// the whole document selector
		return objOrArr, nil
	}
	if !strings.HasPrefix(ptr, "/") {
		return nil, fmt.Errorf("invalid json pointer: %s", ptr)
	}

	tokens := strings.Split(ptr, "/")[1:] // ignore the first empty token
	var root interface{}
	root = objOrArr
	for tIdx, token := range tokens {
		token = Unescape(token)

	typeSwitch:
		if root == nil {
			return nil, fmt.Errorf("value for '%s' not found:\n%s", token, evalMsg(tIdx, tokens))
		}
		switch t := root.(type) {
		case map[string]interface{}:
			if val, ok := t[token]; ok {
				root = val
			} else {
				root = nil
				goto typeSwitch
			}

		case *map[string]interface{}:
			root = *t
			goto typeSwitch
		case []interface{}:
			idx, err := strconv.ParseInt(token, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate %s:%v", ptr, err)
			}
			if idx < 0 || int(idx) >= len(t) {
				return nil, fmt.Errorf("index out of bounds for '%s': %d out of %d", token, idx, len(t))
			}
			root = t[idx]

		case *[]interface{}:
			root = *t
			goto typeSwitch
		}
	}
	return root, nil
}

// evalMsg generates an over-engineered error message
func evalMsg(failedAt int, tokens []JSONPointerToken) string {
	tmp := &strings.Builder{}
	sb := &strings.Builder{}
	for i, t := range tokens {
		sb.WriteString("/")
		sb.WriteString(t)
		if i < failedAt {
			for s := 0; s < len(t); s++ {
				tmp.WriteString(" ")
			}
		} else {
			if i == failedAt {
				tmp.WriteString(" ^")
			}
			for s := 0; s < len(t); s++ {
				tmp.WriteString("~")
			}
		}
	}
	sb.WriteString("\n")
	sb.WriteString(tmp.String())
	return sb.String()
}

// Escapes takes any string and returns a token.
// ~ becomes ~0 and / becomes ~1
func Escape(str string) JSONPointerToken {
	tmp := strings.Replace(str, "~", "~0", -1)
	return strings.Replace(tmp, "/", "~1", -1)
}

// Unescape takes a token and returns the original string.
// ~0 becomes ~ and ~1 becomes /
func Unescape(str JSONPointerToken) string {
	tmp := strings.Replace(str, "~1", "/", -1)
	return strings.Replace(tmp, "~0", "~", -1)
}

// AsString takes a JSONPointer and tries to interpret the result as a string.
// The following rules are applied:
//  * numbers are converted to the according string representation
//  * booleans are converted to true|false
//  * null is converted to the empty string
//  * a non-resolvable value, returns the empty string
//  * arrays and objects are converted into a json string
//  * String() methods are used, if available
//  * Anything else is converted using sprintf and %v directive
func AsString(objOrArr interface{}, ptr JSONPointer) string {
	val, err := Evaluate(objOrArr, ptr)
	if err != nil {
		logger.Info(Fields{"msg": "cannot resolve value", "ptr": ptr})
	}
	return ToString(val)
}

// AsFloat takes a JSONPointer and tries to interpret the result as a float.
// The following rules are applied:
//  * numbers are converted to a float
//  * booleans are converted to 1|0
//  * null is converted NaN
//  * a non-resolvable value, returns also NaN
//  * arrays and objects are converted to NaN
//  * String() method is invoked, if available, and output parsed. Returns NaN if not parsable.
//  * Anything else is converted using sprintf and %v directive and tried to be parsed. Returns NaN if not parsable.
func AsFloat(objOrArr interface{}, ptr JSONPointer) float64 {
	val, err := Evaluate(objOrArr, ptr)
	if err != nil {
		logger.Info(Fields{"msg": "cannot resolve value", "ptr": ptr})
	}
	return ToFloat64(val)
}

// AsArray evaluates the JSONPointer and unwraps either []interface{} or *[]interface{} slices. Any other slice
// types are unboxed to []interface{}.
func AsArray(objOrArr interface{}, ptr JSONPointer) []interface{} {
	val, err := Evaluate(objOrArr, ptr)
	if err != nil {
		logger.Info(Fields{"msg": "cannot resolve value", "ptr": ptr})
	}
	switch t := val.(type) {
	case []interface{}:
		return t
	case *[]interface{}:
		return *t
	default:
		panic("TODO implement a reflect autoboxing")
	}
	return nil
}

// AsFloatArray evaluates the JSONPointer and tries to interpret any slice value as float (see AsFloat) for rules
func AsFloatArray(objOrArr interface{}, ptr JSONPointer) []float64 {
	slice := AsArray(objOrArr, ptr)
	res := make([]float64, len(slice))
	for i, v := range slice {
		res[i] = ToFloat64(v)
	}
	return res
}

// AsStringArray evaluates the JSONPointer and tries to interpret any slice value as string (see AsString) for rules
func AsStringArray(objOrArr interface{}, ptr JSONPointer) []string {
	slice := AsArray(objOrArr, ptr)
	res := make([]string, len(slice))
	for i, v := range slice {
		res[i] = ToString(v)
	}
	return res
}
