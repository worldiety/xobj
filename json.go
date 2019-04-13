package xobj

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type jsonUnmarshaler struct {
}

func (jsonUnmarshaler) Unmarshal(data []byte) (Node, error) {
	if !hasFirstChar(data, '{', 8) {
		return nil, ErrUnsupportedFormat
	}
	res := &jsonNode{asObj: make(map[string]interface{})}
	err := json.Unmarshal(data, &res.asObj)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//==

var _ NodeList = (*jsonList)(nil)

type jsonList []*jsonNode

func (l jsonList) Size() int {
	return len(l)
}

func (l jsonList) Get(idx int) Node {
	return l[idx]
}

//==

var _ Node = (*jsonNode)(nil)

type jsonNode struct {
	name    string
	asObj   map[string]interface{}
	asValue interface{}
	parent  *jsonNode
	nil     bool
}

func (e *jsonNode) IsNull() bool {
	return e == nil || strings.ToLower(e.AsString()) == "null"
}

func (e *jsonNode) Get(name string) Node {
	if e == nil {
		return e.asNil(name)
	}

	// field lookup
	if e.asObj != nil {
		v := e.asObj[name]
		return e.asNode(name, v)
	}

	// array lookup fallback
	if arr, ok := e.asValue.([]interface{}); ok {
		idx, err := strconv.ParseInt(name, 10, 32)
		if err != nil || idx < 0 || int(idx) >= len(arr) {
			return e.asNil(name)
		}
		return e.asNode(name, arr[idx])
	}

	return e.asNil(name)
}

func (e *jsonNode) IsNil() bool {
	return e == nil || e.nil
}

func (e *jsonNode) Remove() bool {
	if e == nil || e.nil == true || e.parent == nil {
		return true
	}
	if e.parent != nil {
		if e.parent.asObj != nil {
			delete(e.parent.asObj, e.name)
			e.parent = nil
			return true
		}
	}
	panic("implementation failure")
}

func (e *jsonNode) AsString() string {
	if e == nil || e.nil {
		return "null"
	}
	if e.asValue == nil {
		return "null"
	}

	if str, ok := e.asValue.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", e.asValue)
}

func (e *jsonNode) SetString(str string) {
	panic("implement me")
}

func (e *jsonNode) Name() string {
	return e.name
}

func (e *jsonNode) asNode(name string, v interface{}) *jsonNode {
	if v == nil {
		return e.asNil(name)
	}
	switch t := v.(type) {
	case map[string]interface{}:
		return &jsonNode{name, t, nil, e, false}
	default:
		return &jsonNode{name, nil, t, e, false}

	}
}

func (e *jsonNode) asNil(name string) *jsonNode {
	return &jsonNode{name, nil, nil, e, true}
}

func (e *jsonNode) Nodes() NodeList {
	res := jsonList{}
	if e == nil || e.nil {
		return res
	}
	if e.asObj != nil {
		for key, val := range e.asObj {
			res = append(res, e.asNode(key, val))
		}
		return res
	}

	if t, ok := e.asValue.([]interface{}); ok {
		for idx, v := range t {
			res = append(res, e.asNode(strconv.Itoa(idx), v))
		}
	}

	return res

}
