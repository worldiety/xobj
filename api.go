package xobj

import "fmt"

var unmarshaler []Unmarshaler
var ErrUnsupportedFormat = fmt.Errorf("unsupported format")

func init() {
	RegisterUnmarshaler(jsonUnmarshaler{})
}

type Unmarshaler interface {
	Unmarshal(data []byte) (Node, error)
}

// RegisterUnmarshaler allows to extend the xobj support for further formats
func RegisterUnmarshaler(u Unmarshaler) {
	unmarshaler = append(unmarshaler, u)
}

// Unmarshal autodetects the kind of format and decodes it into
// a Node
func Unmarshal(data []byte) (Node, error) {
	for _, u := range unmarshaler {
		elem, err := u.Unmarshal(data)
		if err != nil {
			if err != ErrUnsupportedFormat {
				return nil, err
			}

		} else {
			return elem, nil
		}
	}
	return nil, ErrUnsupportedFormat
}

// A Node is more or less the same as the XML ElementNode. It also represents a json field/property.
type Node interface {
	// Name return the name of the node.
	//
	// XML: If node is an Element the elements name is returned. If node is Attr, the attributes name is returned.
	// Any other nodes have no name.
	Name() string

	// Get resolves the first child with the given name.
	//
	// XML: If node is Attr, CDATA, Comment, etc. always a nil node is returned. If node is Element, the first Element
	// with that name is returned, otherwise a nil node. To read an attribute of this node, you need to prefix the name
	// with a :
	//
	// JSON: Simply returns the according field or a nil node. If this node refers to an array and the name
	// looks like suitable index within range, the value is returned.
	Get(name string) Node

	// IsNil checks if this node is a nil node.
	IsNil() bool

	// IsNull returns true if this is a nil node or
	//
	// XML: #AsString() evaluates to "null"
	//
	// JSON: the value is the natural null or #AsString() evaluates to "null"
	//
	// struct: the same as JSON
	IsNull() bool

	// Remove detaches this node from its parent. Returns true if the operation was successful.
	Remove() bool

	// AsString returns a naive string interpretation.
	//
	// XML: If node is an Attr, a Comment or a Text, its content is simply returned. If node is an Element, all
	// Text nodes are concated in order of appearance recursively.
	//
	// JSON: Only returns a naive string interpolation of bool, number and null if it is not an array or object.
	// In the latter, just an empty string is returned.
	AsString() string

	// SetString puts the string, removing any other data.
	//
	// XML: If node is an Attr, it's value is updated. If node is an Element, all Text is removed and replaced
	// by a single text node.
	//
	// JSON: The value is replaced with the given string value
	//
	// struct: It tries to duck type the string into the according field. int64 may become 0, float64 NaN, bool false.
	SetString(str string)

	// Nodes returns a NodeList. Semantic is defined as follows:
	//
	// XML: If node is an Element, it's attributes (prefixed by :) and it's direct child nodes are returned
	//
	// JSON: If node is an Object, it's fields are returned. If node is a primitive (bool, string, number, nil)
	// an empty list is returned. If node is an array, each element in the array is a node in the list.
	//
	// struct: Just the same rules as JSON
	Nodes() NodeList
}

// A NodeList is a collection of arbitrary node types.
type NodeList interface {
	// Size returns the amount of elements in the list
	Size() int

	// Get returns the element at the given index. Panics if out of range.
	Get(idx int) Node
}


