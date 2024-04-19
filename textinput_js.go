//go:build go1.12
// +build go1.12

package goey

import (
	"syscall/js"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/internal/js"
)

type textinputElement struct {
	Control

	onChange   goeyjs.ChangeStringCB
	onFocus    goeyjs.FocusCB
	onBlur     goeyjs.BlurCB
	onEnterKey goeyjs.EnterKeyCB
}

func (w *TextInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := goeyjs.CreateElement("input", "goey form-control")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &textinputElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *textinputElement) Close() {
	w.onChange.Close()
	w.onFocus.Close()
	w.onBlur.Close()
	w.onEnterKey.Close()

	w.Control.Close()
}

func (w *textinputElement) createMeasurementElement() js.Value {
	handle := goeyjs.CreateElement("input", "form-control goey-measure")

	goeyjs.AppendChildToBody(handle)

	return handle
}

func (w *textinputElement) Layout(bc base.Constraints) base.Size {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := base.FromPixelsX(handle.Get("offsetWidth").Int() + 1)
	width = bc.ConstrainWidth(width)
	height := base.FromPixelsY(handle.Get("offsetHeight").Int() + 1)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *textinputElement) MinIntrinsicHeight(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	height := handle.Get("offsetHeight").Int()

	return base.FromPixelsY(height)
}

func (w *textinputElement) MinIntrinsicWidth(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := handle.Get("offsetWidth").Int()

	return base.FromPixelsX(width + 1)
}

func (w *textinputElement) Props() base.Widget {
	return &TextInput{
		Value:       w.handle.Get("value").String(),
		Placeholder: w.handle.Get("placeholder").String(),
		Disabled:    w.handle.Get("disabled").Truthy(),
		Password:    w.handle.Get("type").String() == "password",
		ReadOnly:    w.handle.Get("readonly").Truthy(),
		OnChange:    w.onChange.Fn,
		OnFocus:     w.onFocus.Fn,
		OnBlur:      w.onBlur.Fn,
		OnEnterKey:  w.onEnterKey.Fn,
	}
}

func (w *textinputElement) updateProps(data *TextInput) error {
	w.handle.Set("value", data.Value)
	w.handle.Set("placeholder", data.Placeholder)
	w.handle.Set("disabled", data.Disabled)
	if data.Password {
		w.handle.Set("type", "password")
	} else {
		w.handle.Set("type", "text")
	}
	w.handle.Set("readonly", data.ReadOnly)

	w.onChange.Set(w.handle, data.OnChange)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)
	w.onEnterKey.Set(w.handle, data.OnEnterKey)

	return nil
}
