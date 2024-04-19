package base

var (
	nilKind = NewKind("github.com/chaolihf/goey/base.nil")
)

// Mount will try to mount a widget.  In the case where the widget is non-nil,
// this function is a simple wrapper around calling the method Mount directly.
// If widget is nil, this function will instead return a non-nil element to act
// as a placeholder.
//
// The placeholder element has an intrinsic size of zero and adds no visible
// elements in the GUI.  Unlike other elements, there is no need to call Close,
// as that method is a no-op.
func Mount(parent Control, widget Widget) (Element, error) {
	if widget == nil {
		return (*nilElement)(nil), nil
	}
	return widget.Mount(parent)
}

// MountNil is a wrapper around Mount(parent,nil).
func MountNil() Element {
	return (*nilElement)(nil)
}

type nilElement struct{}

func (*nilElement) Close() {
	// No-op
}

func (*nilElement) Kind() *Kind {
	return &nilKind
}

func (*nilElement) Layout(bc Constraints) Size {
	if bc.IsBounded() {
		return bc.Max
	} else if bc.HasBoundedWidth() {
		return Size{bc.Max.Width, bc.Min.Height}
	} else if bc.HasBoundedHeight() {
		return Size{bc.Min.Width, bc.Max.Height}
	}
	return bc.Min
}

func (*nilElement) MinIntrinsicHeight(Length) Length {
	return 0
}

func (*nilElement) MinIntrinsicWidth(Length) Length {
	return 0
}

func (*nilElement) Props() Widget {
	return nil
}

func (*nilElement) SetBounds(Rectangle) {
	// Do nothing
}

func (*nilElement) UpdateProps(data Widget) error {
	panic("unreachable")
}
