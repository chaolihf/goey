//go:build go1.12
// +build go1.12

package goey

import (
	"strings"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	goeyjs "bitbucket.org/rj/goey/internal/js"
)

type buttonElement struct {
	Control

	onClick goeyjs.ClickCB
	onFocus goeyjs.FocusCB
	onBlur  goeyjs.BlurCB
}

func (w *Button) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "button")
	defer parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &buttonElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *buttonElement) Click() {
	w.handle.Call("click")
}

func (w *buttonElement) Close() {
	w.onClick.Close()
	w.onFocus.Close()
	w.onBlur.Close()

	w.Control.Close()
}

func (w *buttonElement) createMeasurementElement() js.Value {
	document := js.Global().Get("document")

	handle := document.Call("createElement", "button")
	handle.Set("className", "btn btn-primary goey-measure")
	handle.Set("textContent", w.handle.Get("textContent"))

	body := document.Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)

	return handle
}

func (w *buttonElement) Layout(bc base.Constraints) base.Size {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := base.FromPixelsX(handle.Get("offsetWidth").Int() + 1)
	width = bc.ConstrainWidth(width)
	height := base.FromPixelsY(handle.Get("offsetHeight").Int() + 1)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *buttonElement) MinIntrinsicHeight(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	height := handle.Get("offsetHeight").Int()

	return base.FromPixelsY(height)
}

func (w *buttonElement) MinIntrinsicWidth(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := handle.Get("offsetWidth").Int()

	return base.FromPixelsX(width + 1)
}

func (w *buttonElement) Props() base.Widget {
	return &Button{
		Text:     w.handle.Get("textContent").String(),
		Default:  strings.Contains(w.handle.Get("className").String(), "primary"),
		Disabled: w.handle.Get("disabled").Truthy(),
		OnClick:  w.onClick.Fn,
		OnFocus:  w.onFocus.Fn,
		OnBlur:   w.onBlur.Fn,
	}
}

func (w *buttonElement) updateProps(data *Button) error {
	w.handle.Set("textContent", data.Text)
	w.handle.Set("disabled", data.Disabled)
	if data.Default {
		w.handle.Set("className", "goey btn btn-primary")
	} else {
		w.handle.Set("className", "goey btn btn-secondary")
	}
	w.onClick.Set(w.handle, data.OnClick)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
