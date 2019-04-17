package xobj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/worldiety/jsonml"
)

// parsers contain all registered format interpreters
var parsers = make([]Parser, 0)

// A Parser is a drop-in contract to extend the #Parse() method of xobj.
type Parser interface {
	// Parse reads the data and returns an object
	Parse(data []byte) (Obj, error)
}

// RegisterParser accepts more format interpreters for the #Parse() method.
// Registering parsers is not thread safe with #Parse(), so ensure that
// you do that at #init() time.
func RegisterParser(parser Parser) {
	parsers = append(parsers, parser)
}

type parserFunc func(data []byte) (Obj, error)

func (f parserFunc) Parse(data []byte) (Obj, error) {
	return f(data)
}

func init() {
	// parse as json object
	RegisterParser(parserFunc(func(data []byte) (Obj, error) {
		obj := Object{}
		err := json.Unmarshal(data, &obj)
		return obj, err
	}))

	// parse as json array
	RegisterParser(parserFunc(func(data []byte) (Obj, error) {
		obj := Object{}
		arr := &Array{}
		obj.Put("array", arr)

		err := json.Unmarshal(data, &arr)
		return obj, err
	}))

	// parse as xml (jsonml transformation)
	RegisterParser(parserFunc(func(data []byte) (Obj, error) {
		obj := Object{}
		slice, err := jsonml.ToJSON(true, bytes.NewReader(data))
		if slice != nil {
			obj.Put("xml", slice)
		}
		return obj, err
	}))

}

// Parse tries to parse the given bytes either as XML or as JSON.
// If data looks like XML, a jsonml transformation is applied, which
// is available in the field 'xml'.
// If data represents an array, it is wrapped automatically into
// an object, using the field name "array".
//
// You can extend the capabilities by registering your custom interpreter using #RegisterParser()
func Parse(data []byte) (Obj, error) {
	for _, p := range parsers {
		obj, err := p.Parse(data)
		if err == nil {
			return obj, nil
		}
	}
	return nil, fmt.Errorf("xobj: unsupported format")
}
