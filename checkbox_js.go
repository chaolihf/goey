// +build go1.12

package goey

import (
	"fmt"
	"math/rand"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/js"
)

type checkboxElement struct {
	Control
	elemInput js.Value
	elemLabel js.Value

	onChange goeyjs.ChangeBoolCB
	onFocus  goeyjs.FocusCB
	onBlur   goeyjs.BlurCB
}

func (w *Checkbox) mount(parent base.Control) (base.Element, error) {
	id := fmt.Sprintf("goey%x", rand.Uint64())

	// Create the control
	handle := js.Global().Get("document").Call("createElement", "div")
	handle.Set("className", "goey form-check")
	elemInput := js.Global().Get("document").Call("createElement", "input")
	elemInput.Set("type", "checkbox")
	elemInput.Set("className", "form-check-input")
	elemInput.Set("id", id)
	handle.Call("appendChild", elemInput)
	elemLabel := js.Global().Get("document").Call("createElement", "label")
	elemLabel.Set("className", "form-check-label")
	elemLabel.Set("htmlFor", id)
	handle.Call("appendChild", elemLabel)
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &checkboxElement{
		Control:   Control{handle},
		elemInput: elemInput,
		elemLabel: elemLabel,
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *checkboxElement) Click() {
	w.elemInput.Call("click")
}

func (w *checkboxElement) Close() {
	w.onChange.Close()
	w.onFocus.Close()
	w.onBlur.Close()

	w.Control.Close()
}

func (w *checkboxElement) createMeasurementElement() js.Value {
	document := js.Global().Get("document")

	handle := document.Call("createElement", "div")
	handle.Set("className", "form-check")
	elemInput := document.Call("createElement", "input")
	elemInput.Set("type", "checkbox")
	elemInput.Set("className", "form-check-input")
	handle.Call("appendChild", elemInput)
	elemLabel := document.Call("createElement", "label")
	elemLabel.Set("className", "form-check-label")
	handle.Call("appendChild", elemLabel)
	handle.Get("style").Set("visibility", "hidden")

	body := document.Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)

	return handle
}

func (w *checkboxElement) Layout(bc base.Constraints) base.Size {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := base.FromPixelsX(handle.Get("offsetWidth").Int() + 1)
	width = bc.ConstrainWidth(width)
	height := base.FromPixelsY(handle.Get("offsetHeight").Int() + 1)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *checkboxElement) MinIntrinsicHeight(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	height := handle.Get("offsetHeight").Int()

	return base.FromPixelsY(height)
}

func (w *checkboxElement) MinIntrinsicWidth(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := handle.Get("offsetWidth").Int()

	return base.FromPixelsX(width + 1)
}

func (w *checkboxElement) Props() base.Widget {
	return &Checkbox{
		Text:     w.elemLabel.Get("textContent").String(),
		Value:    w.elemInput.Get("checked").Truthy(),
		Disabled: w.elemInput.Get("disabled").Truthy(),
		OnChange: w.onChange.Fn,
		OnFocus:  w.onFocus.Fn,
		OnBlur:   w.onBlur.Fn,
	}
}

func (w *checkboxElement) TakeFocus() bool {
	w.elemInput.Call("focus")
	return true
}

func (w *checkboxElement) updateProps(data *Checkbox) error {
	w.elemInput.Set("checked", data.Value)
	w.elemInput.Set("disabled", data.Disabled)
	w.elemLabel.Set("textContent", data.Text)
	w.onChange.Set(w.elemInput, data.OnChange)
	w.onFocus.Set(w.elemInput, data.OnFocus)
	w.onBlur.Set(w.elemInput, data.OnBlur)

	return nil
}
