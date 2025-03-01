package goey

import (
	"errors"
	"testing"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/mock"
)

func (w *PaddingElement) Props() base.Widget {
	child := base.Widget(nil)
	if w.child != nil {
		child = w.child.(Proper).Props()
	}

	return &Padding{
		Insets: w.insets,
		Child:  child,
	}
}

func TestPaddingMount(t *testing.T) {
	// These should all be able to mount without error.
	testMountWidgets(t,
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	)

	// These should mount with an error.
	err := errors.New("Mock error 1")
	testMountWidgetsFail(t, err,
		&Padding{Child: &mock.Widget{Err: err}},
	)
	testMountWidgetsFail(t, err,
		&Padding{Insets: DefaultInsets(), Child: &mock.Widget{Err: err}},
	)
}

func TestPaddingClose(t *testing.T) {
	testCloseWidgets(t,
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	)
}

func TestPaddingUpdateProps(t *testing.T) {
	testUpdateWidgets(t, []base.Widget{
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	}, []base.Widget{
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "AB"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "BC"}},
		&Padding{},
		&Padding{Child: &Button{Text: "CD"}},
	})
}

func TestPaddingMinIntrinsicSize(t *testing.T) {
	size1 := base.Size{10 * DIP, 20 * DIP}
	size2 := base.Size{15 * DIP, 25 * DIP}
	sizeZ := base.Size{}
	insets := Insets{1 * DIP, 2 * DIP, 3 * DIP, 4 * DIP}

	cases := []struct {
		mockSize base.Size
		insets   Insets
		out      base.Size
	}{
		{size1, Insets{}, size1},
		{size2, Insets{}, size2},
		{sizeZ, Insets{}, sizeZ},
		{size1, insets, base.Size{16 * DIP, 24 * DIP}},
		{size2, insets, base.Size{21 * DIP, 29 * DIP}},
		{sizeZ, insets, base.Size{6 * DIP, 4 * DIP}},
	}

	for i, v := range cases {
		elem := PaddingElement{
			child:  mock.NewIfNotZero(v.mockSize),
			insets: v.insets,
		}

		if out := elem.MinIntrinsicWidth(base.Inf); out != v.out.Width {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Width)
		}
		if out := elem.MinIntrinsicHeight(base.Inf); out != v.out.Height {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Height)
		}
	}
}
