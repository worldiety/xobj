package xobj

import (
	"testing"
)

func TestQuery(t *testing.T) {
	obj := NewObj()
	obj.OptArray("persons").
		AddObject(NewObj().Put("id", 1).Put("name", "Frodo")).
		AddObject(NewObj().Put("id", 2).Put("name", "Sam")).
		AddObject(NewObj().Put("id", 3).Put("name", "Gandalf"))

	// get all persons
	if res, err := Query(obj, "$.persons.*"); err != nil || res.Size() != 3 {
		t.Fatal("unexpected", res, err)
	}

	// get by name
	if res, err := Query(obj, `$..[?(@.name=="Sam")]`); err != nil || res.Size() != 1 {
		t.Fatal("unexpected", res, err)
	}

	// get by id => TODO does not work: the library is broken
	//if res, err := Query(obj, `$..[?(@.id==1)]`); err != nil || res.Size() != 1 {
	//	t.Fatal("unexpected", res, err)
	//}

	// get by id smaller than one => TODO does not work: the library is broken
	//if res, err := Query(obj, `$..[?(@.id<1)]`); err != nil || res.Size() != 1 {
	//	t.Fatal("unexpected", res, err)
	//}

	if res, err := Query(obj, `$..persons[1]`); err != nil || res.Size() != 1 {
		t.Fatal("unexpected", res, err)
	}

	if res, err := Query(obj, `$.persons..name`); err != nil || res.Size() != 3 {
		t.Fatal("unexpected", res, err)
	}

	if res, err := Query(obj, `$..name`); err != nil || res.Size() != 3 {
		t.Fatal("unexpected", res, err)
	}

}
