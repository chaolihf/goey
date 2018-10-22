// +build !gnustep

package goey

import (
	"errors"
	"image/color"
	"testing"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/mock"
)

var (
	black = color.RGBA{0, 0, 0, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
	red   = color.RGBA{0xcc, 0xaa, 0x88, 0xff}
)

func (w *decorationElement) Props() base.Widget {
	widget := w.props()
	if w.child != nil {
		widget.Child = w.child.(Proper).Props()
	}

	return widget
}

func decorationChildWidget(child base.Element) base.Widget {
	if child == nil {
		return nil
	}

	return child.(Proper).Props()
}

func TestDecorationCreate(t *testing.T) {
	// These should all be able to mount without error.
	testingRenderWidgets(t,
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
		&Decoration{Fill: white, Stroke: black},
		&Decoration{Fill: red},
	)

	// These should mount with an error.
	err := errors.New("Mock error 1")
	testingRenderWidgetsFail(t, err,
		&Decoration{Child: &mock.Widget{Err: err}},
	)
	testingRenderWidgetsFail(t, err,
		&Decoration{Insets: DefaultInsets(), Child: &mock.Widget{Err: err}},
	)
}

func TestDecorationClose(t *testing.T) {
	testingCloseWidgets(t,
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
	)
}

func TestDecorationUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
		&Decoration{Fill: white, Stroke: black},
		&Decoration{Fill: red},
	}, []base.Widget{
		&Decoration{},
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: black},
		&Decoration{Fill: white},
	})
}

func TestDecorationMinIntrinsicSize(t *testing.T) {
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
		elem := decorationElement{
			insets: v.insets,
		}
		if !v.mockSize.IsZero() {
			elem.child = mock.New(v.mockSize)
		}

		if out := elem.MinIntrinsicWidth(base.Inf); out != v.out.Width {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Width)
		}
		if out := elem.MinIntrinsicHeight(base.Inf); out != v.out.Height {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Height)
		}
	}
}
