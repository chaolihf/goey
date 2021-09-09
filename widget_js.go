package goey

import (
	"fmt"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"gitlab.com/stone.code/assert"
)

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkWidget when targeting GTK.
//
// Unless developping new widgets, users should not need to use this type.
//
// Any method's on this type will be platform specific.
type Control struct {
	handle js.Value
}

// Close removes the element from the GUI, and frees any associated resources.
func (w *Control) Close() {
	if !w.handle.IsNull() {
		w.handle.Call("remove")
		w.handle = js.Null()
	}
}

// Layout determines the best size for an element that satisfies the
// constraints.
func (w *Control) Layout(bc base.Constraints) base.Size {
	if !bc.HasBoundedWidth() && !bc.HasBoundedHeight() {
		// No need to worry about breaking the constraints.  We can take as
		// much space as desired.
		width, height := 100, 16
		// Dimensions may need to be increased to meet minimums.
		return bc.Constrain(base.Size{base.FromPixelsX(width), base.FromPixelsY(height)})
	}
	if !bc.HasBoundedHeight() {
		// No need to worry about height.  Find the width that best meets the
		// widgets preferred width.
		width1 := 100 // gtk.WidgetNaturalWidth(w.handle)
		width := bc.ConstrainWidth(base.FromPixelsX(width1))
		// Get the best height for this width.
		height := 16 // gtk.WidgetNaturalHeightForWidth(w.handle, width.PixelsX())
		// Height may need to be increased to meet minimum.
		return base.Size{width, bc.ConstrainHeight(base.FromPixelsY(height))}
	}

	// Not clear the following is the best general approach given GTK layout
	// model.
	width, height := 100, 16 // gtk.WidgetNaturalSize(w.handle)
	return bc.Constrain(base.Size{base.FromPixelsX(width), base.FromPixelsY(height)})
}

// MinIntrinsicHeight returns the minimum height that this element requires
// to be correctly displayed.
func (w *Control) MinIntrinsicHeight(width base.Length) base.Length {
	if width != base.Inf {
		height := 16 // gtk.WidgetMinHeightForWidth(w.handle, width.PixelsX())
		return base.FromPixelsY(height)
	}
	height := 16 // gtk.WidgetMinHeight(w.handle)
	return base.FromPixelsY(height)
}

// MinIntrinsicWidth returns the minimum width that this element requires
// to be correctly displayed.
func (w *Control) MinIntrinsicWidth(base.Length) base.Length {
	width := 100 // gtk.WidgetMinWidth(w.handle)
	return base.FromPixelsX(width)
}

// SetBounds updates the position of the widget.
func (w *Control) SetBounds(bounds base.Rectangle) {
	pixels := bounds.Pixels()
	assert.Assert(pixels.Dx() > 0 && pixels.Dy() > 0, "zero width or zero height bounds for control")

	style := w.handle.Get("style")

	style.Set("left", fmt.Sprintf("%dpx", pixels.Min.X))
	style.Set("top", fmt.Sprintf("%dpx", pixels.Min.Y))
	style.Set("width", fmt.Sprintf("%dpx", pixels.Dx()))
	style.Set("height", fmt.Sprintf("%dpx", pixels.Dy()))
}

func (w *Control) TakeFocus() bool {
	w.handle.Call("focus")
	return true
}
