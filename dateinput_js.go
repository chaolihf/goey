// +build go1.12

package goey

import (
	"syscall/js"
	"time"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/js"
)

type dateinputElement struct {
	Control

	onChange goeyjs.ChangeDateCB
	onFocus  goeyjs.FocusCB
	onBlur   goeyjs.BlurCB
}

func (w *DateInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "input")
	handle.Get("style").Set("position", "absolute")
	handle.Set("type", "date")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &dateinputElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *dateinputElement) Props() base.Widget {
	value, _ := time.Parse("2006-1-2", w.handle.Get("value").String())

	return &DateInput{
		Value:    value.Local(),
		Disabled: w.handle.Get("disabled").Truthy(),
		OnChange: w.onChange.Fn,
		OnFocus:  w.onFocus.Fn,
		OnBlur:   w.onBlur.Fn,
	}
}

func (w *dateinputElement) updateProps(data *DateInput) error {
	w.handle.Set("value", data.Value.Format("2006-01-02"))
	w.handle.Set("disabled", data.Disabled)

	w.onChange.Set(w.handle, data.OnChange)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
