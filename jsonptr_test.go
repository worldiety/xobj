package xobj

import (
	"encoding/json"
	"reflect"
	"testing"
)

const rfcTestJson = `

   {
      "foo": ["bar", "baz"],
      "": 0,
      "a/b": 1,
      "c%d": 2,
      "e^f": 3,
      "g|h": 4,
      "i\\j": 5,
      "k\"l": 6,
      " ": 7,
      "m~n": 8
   }
`

func TestEvaluate(t *testing.T) {
	obj := make(map[string]interface{})
	err := json.Unmarshal([]byte(rfcTestJson), &obj)
	if err != nil {
		t.Fatal(err)
	}
	expectStr(t, obj, "/foo/0", "bar")
	expectFloat64(t, obj, "/", 0)
	expectFloat64(t, obj, "/a~1b", 1)
	expectFloat64(t, obj, "/c%d", 2)
	expectFloat64(t, obj, "/e^f", 3)
	expectFloat64(t, obj, "/g|h", 4)
	expectFloat64(t, obj, "/i\\j", 5)
	expectFloat64(t, obj, "/k\"l", 6)
	expectFloat64(t, obj, "/ ", 7)
	expectFloat64(t, obj, "/m~0n", 8)

	res, err := Evaluate(obj, "")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(obj, res) {
		t.Fatal("expected the same")
	}

	res, err = Evaluate(obj, "/foo")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, ([]interface{}{"bar", "baz"})) {
		t.Fatal("expected the same")
	}

	res, err = Evaluate(obj, "/abc/asd")
	if err != nil {
		t.Fatal(err)
	}
}

func expectStr(t *testing.T, json interface{}, ptr JSONPointer, val string) {
	t.Helper()
	v, err := Evaluate(json, ptr)
	if err != nil {
		t.Fatal(err)
	}

	if str, ok := v.(string); ok {
		if str != val {
			t.Fatal("expected", val, "but got", str)
		}
	} else {
		t.Fatal("unexpected", v)
	}
}

func expectFloat64(t *testing.T, json interface{}, ptr JSONPointer, val float64) {
	t.Helper()
	v, err := Evaluate(json, ptr)
	if err != nil {
		t.Fatal(err)
	}

	if str, ok := v.(float64); ok {
		if str != val {
			t.Fatal("expected", val, "but got", str)
		}
	} else {
		t.Fatal("unexpected", v)
	}
}
