package goey

import (
	"fmt"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
)

type progressElement struct {
	Control

	innerDiv js.Value
}

func (w *Progress) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "div")
	handle.Set("className", "progress")
	handle.Get("style").Set("position", "absolute")
	handle2 := js.Global().Get("document").Call("createElement", "div")
	handle2.Set("className", "progress-bar")
	handle2.Set("role", "progress-bar")
	handle.Call("appendChild", handle2)
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &progressElement{
		Control:  Control{handle},
		innerDiv: handle2,
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *progressElement) Props() base.Widget {
	return &Progress{
		Value: w.innerDiv.Get("aria-valuenow").Int(),
		Min:   w.innerDiv.Get("aria-valuemin").Int(),
		Max:   w.innerDiv.Get("aria-valuemax").Int(),
	}
}

func (w *progressElement) updateProps(data *Progress) error {
	w.innerDiv.Get("style").Set("width", fmt.Sprintf("%f%%", float32(data.Value-data.Min)/float32(data.Max-data.Min)))
	w.innerDiv.Set("aria-valuenow", data.Value)
	w.innerDiv.Set("aria-valuemin", data.Min)
	w.innerDiv.Set("aria-valuemax", data.Max)
	return nil
}
