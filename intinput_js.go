//go:build go1.12
// +build go1.12

package goey

import (
	"strconv"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/js"
)

type intinputElement struct {
	Control

	onChange   goeyjs.ChangeInt64CB
	onFocus    goeyjs.FocusCB
	onBlur     goeyjs.BlurCB
	onEnterKey goeyjs.EnterKeyInt64CB
}

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "input")
	handle.Set("className", "goey form-control")
	handle.Set("type", "number")
	handle.Set("step", 1)
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &intinputElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *intinputElement) Close() {
	w.onChange.Close()
	w.onFocus.Close()
	w.onBlur.Close()
	w.onEnterKey.Close()

	w.Control.Close()
}

func (w *intinputElement) createMeasurementElement() js.Value {
	document := js.Global().Get("document")

	handle := document.Call("createElement", "input")
	handle.Set("className", "form-control goey-measure")
	handle.Set("type", "number")
	handle.Set("step", 1)

	body := document.Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)

	return handle
}

func (w *intinputElement) Layout(bc base.Constraints) base.Size {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := base.FromPixelsX(handle.Get("offsetWidth").Int() + 1)
	width = bc.ConstrainWidth(width)
	height := base.FromPixelsY(handle.Get("offsetHeight").Int() + 1)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *intinputElement) MinIntrinsicHeight(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	height := handle.Get("offsetHeight").Int()

	return base.FromPixelsY(height)
}

func (w *intinputElement) MinIntrinsicWidth(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := handle.Get("offsetWidth").Int()

	return base.FromPixelsX(width + 1)
}

func (w *intinputElement) Props() base.Widget {
	value, _ := strconv.ParseInt(
		w.handle.Get("value").String(),
		10, 64)
	min, _ := strconv.ParseInt(
		w.handle.Get("min").String(),
		10, 64)
	max, _ := strconv.ParseInt(
		w.handle.Get("max").String(),
		10, 64)

	return &IntInput{
		Value:       value,
		Min:         min,
		Max:         max,
		Placeholder: w.handle.Get("placeholder").String(),
		Disabled:    w.handle.Get("disabled").Truthy(),
		OnChange:    w.onChange.Fn,
		OnFocus:     w.onFocus.Fn,
		OnBlur:      w.onBlur.Fn,
		OnEnterKey:  w.onEnterKey.Fn,
	}
}

func (w *intinputElement) updateProps(data *IntInput) error {
	w.handle.Set("value", data.Value)
	w.handle.Set("placeholder", data.Placeholder)
	w.handle.Set("disabled", data.Disabled)
	w.handle.Set("min", data.Min)
	w.handle.Set("max", data.Max)

	w.onChange.Set(w.handle, data.OnChange)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
