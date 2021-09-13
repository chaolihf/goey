// +build go1.12

package goey

import (
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	goeyjs "bitbucket.org/rj/goey/internal/js"
)

type textareaElement struct {
	Control

	minLines int
	onChange goeyjs.ChangeStringCB
	onFocus  goeyjs.FocusCB
	onBlur   goeyjs.BlurCB
}

func (w *TextArea) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "textarea")
	handle.Set("className", "goey form-control")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &textareaElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *textareaElement) Close() {
	w.onChange.Close()
	w.onFocus.Close()
	w.onBlur.Close()

	w.Control.Close()
}

func (w *textareaElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(bc.Max.Width)
	width = bc.ConstrainWidth(width)
	height := w.MinIntrinsicHeight(width)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *textareaElement) MinIntrinsicHeight(base.Length) base.Length {
	// Create a dummy button
	handle := js.Global().Get("document").Call("createElement", "textarea")
	handle.Set("value", w.handle.Get("value"))
	handle.Get("style").Set("visibility", "hidden")

	body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)
	height := handle.Get("offsetHeight").Int()
	handle.Call("remove")

	return base.FromPixelsY(height)
}

func (w *textareaElement) MinIntrinsicWidth(base.Length) base.Length {
	// Create a dummy button
	handle := js.Global().Get("document").Call("createElement", "button")
	handle.Set("className", "btn btn-primary")
	handle.Set("value", w.handle.Get("value"))
	handle.Get("style").Set("visibility", "hidden")

	body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)
	width := handle.Get("offsetWidth").Int()
	handle.Call("remove")

	return base.FromPixelsX(width + 1)
}

func (w *textareaElement) Props() base.Widget {
	return &TextArea{
		Value:       w.handle.Get("value").String(),
		Placeholder: w.handle.Get("placeholder").String(),
		Disabled:    w.handle.Get("disabled").Truthy(),
		MinLines:    w.minLines,
		OnChange:    w.onChange.Fn,
		OnFocus:     w.onFocus.Fn,
		OnBlur:      w.onBlur.Fn,
	}
}

func (w *textareaElement) updateProps(data *TextArea) error {
	w.handle.Set("value", data.Value)
	w.handle.Set("disabled", data.Disabled)
	w.handle.Set("placeholder", data.Placeholder)
	w.minLines = data.MinLines
	w.onChange.Set(w.handle, data.OnChange)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
