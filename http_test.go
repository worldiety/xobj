package xobj

import (
	"testing"
	"time"
)

func TestRequestBuilder_Get(t *testing.T) {
	obj, status, err := NewRequest().Host("test.worldiety.org").Path("test.xml").Get()
	if err != nil {
		t.Fatal(err)
	}
	if status != 200 {
		t.Fatal("unexpected", status)
	}
	if !obj.Has("xml") {
		t.Fatal("expected xml")
	}
}

func TestRequestBuilder_Get2(t *testing.T) {
	_, status, err := NewRequest().
		Host("test.worldiety.org").
		Path("sleep.php").
		ResponseHeaderTimeout(3 * time.Second).
		Get()

	if status != 0 {
		t.Fatal("unexpected", status)
	}

	if err == nil {
		t.Fatal("expected timeout")
	}
}

func TestRequestBuilder_Get3(t *testing.T) {
	_, status, err := NewRequest().
		Host("test.worldiety.org").
		Path("sleep.php").
		DialTimeout(3 * time.Second).
		Get()

	if status != 0 {
		t.Fatal("unexpected", status)
	}

	if err == nil {
		t.Fatal("expected timeout")
	}
}

func TestRequestBuilder_Get4(t *testing.T) {
	_, status, _ := NewRequest().
		Host("test.worldiety.org").
		Path("not found").
		Get()

	if status != 404 {
		t.Fatal("unexpected", status)
	}

}
