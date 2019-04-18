package xobj

import "github.com/PaesslerAG/jsonpath"

// Query parses the query string as jsonPath (see https://goessner.net/articles/JsonPath/)
// and applies it on the given object.
//
// This is currently delegated to https://github.com/PaesslerAG/jsonpath and requires quite
// expensive transformations.
func Query(obj Obj, query string) (Arr, error) {
	res, err := jsonpath.Get(query, UnwrapObj(obj))
	if err != nil {
		return nil, err
	}
	slice := res.([]interface{})
	tmp := Array(slice)
	return &tmp, nil
}
