package base

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func ExampleMount() {
	// This won't work in real code, as the zero value for a control is not
	// generally useable.
	parent := Control{}

	widget := &mock{}
	elem, err := Mount(parent, widget)
	if err != nil {
		panic("Unexpected error!")
	}
	defer elem.Close()

	if widget.Kind() != elem.Kind() {
		panic("Internal error, kinds do not match")
	}

	fmt.Println("OK")

	// Output:
	// OK
}

func ExampleMount_nil() {
	// This won't work in real code, as the zero value for a control is not
	// generally useable.
	parent := Control{}

	// It is okay to mount a nil widget.
	elem, err := Mount(parent, nil)
	if err != nil {
		panic("Unexpected error!")
	}
	defer elem.Close()
	fmt.Println("The value of elem is nil...", elem == nil)
	fmt.Println("The kind of elem is...", elem.Kind())
	fmt.Println("The minimum intrinsic height is...", elem.MinIntrinsicHeight(Inf))
	fmt.Println("The minimum intrinsic width is...", elem.MinIntrinsicWidth(Inf))

	// Output:
	// The value of elem is nil... false
	// The kind of elem is... github.com/chaolihf/goey/base.nil
	// The minimum intrinsic height is... 0:00
	// The minimum intrinsic width is... 0:00
}

func TestMount(t *testing.T) {
	kind1 := NewKind("github.com/chaolihf/goey/base.Mock1")
	kind2 := NewKind("github.com/chaolihf/goey/base.Mock2")
	err1 := errors.New("fake error 1 for mounting widget")
	err2 := errors.New("fake error 2 for mounting widget")

	cases := []struct {
		in  Widget
		out Element
		err error
	}{
		{nil, (*nilElement)(nil), nil},
		{&mock{kind: &kind1}, &mockElement{kind: &kind1}, nil},
		{&mock{kind: &kind1, Prop: 3}, &mockElement{kind: &kind1, Prop: 3}, nil},
		{&mock{kind: &kind2}, &mockElement{kind: &kind2}, nil},
		{&mock{kind: &kind2, Prop: 13}, &mockElement{kind: &kind2, Prop: 13}, nil},
		{&mock{kind: &kind1, err: err1}, nil, err1},
		{&mock{kind: &kind1, err: err2}, nil, err2},
	}

	for i, v := range cases {
		out, err := Mount(Control{}, v.in)
		if err != v.err {
			t.Errorf("Case %d: Returned error does not match, got %v, want %v", i, err, v.err)
		}
		if !reflect.DeepEqual(out, v.out) {
			t.Errorf("Case %d: Returned element does not match, got %v, want %v", i, out, v.out)
		}
	}
}

func TestMountNil(t *testing.T) {
	out, err := Mount(Control{}, nil)
	if err != nil {
		t.Fatalf("Failed to mount <nil> widget: %s", err)
	}
	if out == nil {
		t.Fatalf("Element is nil")
	}

	props := out.(interface{ Props() Widget }).Props()
	if props != nil {
		t.Errorf("Nil element unexpectedly return nil for props")
	}
}
