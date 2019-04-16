package xobj

import (
	"fmt"
	"testing"
)

func TestNewObj(t *testing.T) {
	tmp := NewObj().
		SetBool("flag", true).
		SetInt64("int", 123).
		SetFloat64("float", 3.14).
		SetString("text", "hello world").
		SetObject("obj", NewObj().Put("say", "hello")).
		SetArray("arr", NewArr().
			AddBool(true).
			AddBool(false).
			AddFloat64(6.28).
			AddInt64(246).
			AddString("text in array").
			AddArray(NewArr().AddInt64(1).AddInt64(2).AddInt64(3)).
			AddObject(NewObj().SetString("object in", "array")))
	fmt.Println(tmp.String())
}

const json0 = `
{
	"hello":"world",
	"list":[
		"so",
		"what",
		2,
		3.14,
		true,
		false,
		{"k":"v"},
		[1,2,3]
	]
}
`

const json1 = `
[
	"so",
	"what",
	2,
	3.14,
	true,
	false,
	{"k":"v"},
	[1,2,3]
]
`

const xml0 = `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
	<title>This is an example</title>
	
	<details>
		Something
	
		more
		with
		
		a lot
		of breaks!
	</details>
	<!-- this is an xml comment with < and > and ]] and [[ -->
	
	
	<table caption="a tablet with fruits">
	  <tr>
		<td>0a</td>
		<td>0b</td>
	  </tr>
       <tr>
		<td>1a</td>
		<td>1b</td>
	  </tr>
	</table>
	
	<table>
	  <name>A table desk</name>
	  <width>60</width>
	  <length>113</length>
	</table>
</root>
`

func TestParse(t *testing.T) {
	obj, err := Parse([]byte(json0))
	if err != nil {
		t.Fatal(err)
	}

	list := obj.Keys()
	if list.Size() != 2 {
		t.Fatal("unexpected", list.Size())
	}

	for i := 0; i < list.Size(); i++ {
		if list.Get(i) != "hello" && list.Get(i) != "list" {
			t.Fatal("unexpected", list.Get(i))
		}
	}

	arr, err := obj.AsArray("list")
	if err != nil {
		t.Fatal(err)
	}

	if v, err := arr.AsString(0); err != nil || v != "so" {
		t.Fatal("unexpected", v, err)
	}

	if v, err := arr.AsString(1); err != nil || v != "what" {
		t.Fatal("unexpected", v, err)
	}

	if v, err := arr.AsInt64(2); err != nil || v != 2 {
		t.Fatal("unexpected", v, err)
	}

	if v, err := arr.AsFloat64(3); err != nil || v != 3.14 {
		t.Fatal("unexpected", v, err)
	}

	if v, err := arr.AsBool(4); err != nil || v != true {
		t.Fatal("unexpected", v, err)
	}

	if v, err := arr.AsBool(5); err != nil || v != false {
		t.Fatal("unexpected", v, err)
	}

	if v, err := arr.AsObject(6); err != nil || v != nil {
		if v != nil {
			if str, err := v.AsString("k"); err != nil || str != "v" {
				t.Fatal("unexpected", str, err)
			}
		} else {
			t.Fatal("unexpected", v, err)
		}

	}

	if v, err := arr.AsArray(7); err != nil || v != nil {
		if v != nil {
			if i, err := v.AsInt64(0); err != nil || i != 1 {
				t.Fatal("unexpected", i, err)
			}
			if i, err := v.AsFloat64(1); err != nil || i != 2 {
				t.Fatal("unexpected", i, err)
			}

			if i, err := v.AsString(2); err != nil || i != "3" {
				t.Fatal("unexpected", i, err)
			}

		} else {
			t.Fatal("unexpected", v, err)
		}

	}

	fmt.Println(obj.String())
}

func TestParse2(t *testing.T) {
	obj, err := Parse([]byte(json1))
	if err != nil {
		t.Fatal(err)
	}
	arr, err := obj.AsArray("array")
	if err != nil {
		t.Fatal(err)
	}
	if i, err := arr.AsInt64(2); err != nil || i != 2 {
		t.Fatal("unexpected", i, err)
	}
}

func TestParse3(t *testing.T) {
	obj, err := Parse([]byte(xml0))
	if err != nil {
		t.Fatal(err)
	}
	arr, err := obj.AsArray("xml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(arr.String())
}
