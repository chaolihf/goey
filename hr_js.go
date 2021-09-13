// +build go1.12

package goey

import (
	"fmt"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
)

type hrElement struct {
	Control
}

func (w *HR) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "hr")
	handle.Set("className", "goey")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &hrElement{
		Control: Control{handle},
	}

	return retval, nil
}

func (w *hrElement) SetBounds(bounds base.Rectangle) {
	pixels := bounds.Pixels()

	top := (pixels.Min.Y + pixels.Max.Y) / 2

	style := w.handle.Get("style")

	style.Set("left", fmt.Sprintf("%dpx", pixels.Min.X))
	style.Set("top", fmt.Sprintf("%dpx", top))
	style.Set("width", fmt.Sprintf("%dpx", pixels.Dx()))
	style.Set("height", fmt.Sprintf("1px"))
}
