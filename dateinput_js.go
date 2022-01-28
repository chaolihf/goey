//go:build go1.12
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
	handle := goeyjs.CreateElement("input", "goey form-control")
	handle.Set("type", "date")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &dateinputElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *dateinputElement) Close() {
	w.onChange.Close()
	w.onFocus.Close()
	w.onBlur.Close()

	w.Control.Close()
}

func (w *dateinputElement) createMeasurementElement() js.Value {
	handle := goeyjs.CreateElement("input", "form-control goey-measure")
	handle.Set("type", "date")

	goeyjs.AppendChildToBody(handle)

	return handle
}

func (w *dateinputElement) Layout(bc base.Constraints) base.Size {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := base.FromPixelsX(handle.Get("offsetWidth").Int() + 1)
	width = bc.ConstrainWidth(width)
	height := base.FromPixelsY(handle.Get("offsetHeight").Int() + 1)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *dateinputElement) MinIntrinsicHeight(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	height := handle.Get("offsetHeight").Int()

	return base.FromPixelsY(height)
}

func (w *dateinputElement) MinIntrinsicWidth(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := handle.Get("offsetWidth").Int()

	return base.FromPixelsX(width + 1)
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
