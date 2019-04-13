package xobj

import (
	"testing"
)

var json0 = `
{
  ".comment": "a comment",
  "name": "Frodo",
  "name.firstname": true,
  "age": 42,
  "address": {
    "city": "Hobbingen"
  },
  "hobbies": ["eating","sleeping","reading"],
  "missing": null
}
`

func Test_jsonUnmarsaler_Unmarshal(t *testing.T) {
	node, err := Unmarshal([]byte(json0))
	if err != nil {
		t.Fatal(err)
	}

	type tupel struct {
		value  string
		getter func() string
	}

	set := []tupel{
		{"a comment", func() string {
			return node.Get(".comment").AsString()
		}},

		{"Frodo", func() string {
			return node.Get("name").AsString()
		}},

		{"true", func() string {
			return node.Get("name.firstname").AsString()
		}},

		{"42", func() string {
			return node.Get("age").AsString()
		}},

		{"Hobbingen", func() string {
			return node.Get("address").Get("city").AsString()
		}},
		{"null", func() string {
			return node.Get("missing").AsString()
		}},
	}

	for _, x := range set {
		got := x.getter()
		if got != x.value {
			t.Fatal("expected", x.value, "but got", got)
		}
	}

	nl := node.Get("hobbies").Nodes()
	if nl.Size() != 3 {
		t.Fatal("got", nl.Size())
	}
	for i := 0; i < nl.Size(); i++ {
		switch i {
		case 0:
			if nl.Get(i).AsString() != "eating" {
				t.Fatal("got", nl.Get(i).AsString())
			}
		case 1:
			if nl.Get(i).AsString() != "sleeping" {
				t.Fatal("got", nl.Get(i).AsString())
			}
		case 2:
			if nl.Get(i).AsString() != "reading" {
				t.Fatal("got", nl.Get(i).AsString())
			}
		}
	}

	if node.Get("hobbies").Get("0").AsString() != "eating" {
		t.Fatal("got", node.Get("hobbies").Get("0").AsString())
	}

	nl = node.Nodes()
	if nl.Size() != 7 {
		t.Fatal("got", nl.Size())
	}
}
