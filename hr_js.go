// +build go1.12

package goey

import (
	"syscall/js"

	"bitbucket.org/rj/goey/base"
)

type hrElement struct {
	Control
}

func (w *HR) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "hr")
	handle.Get("style").Set("position", "absolute")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &hrElement{
		Control: Control{handle},
	}

	return retval, nil
}
