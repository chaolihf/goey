//go:build go1.12
// +build go1.12

package goey

import (
	"syscall/js"

	"github.com/chaolihf/goey/base"
	goeyjs "github.com/chaolihf/goey/internal/js"
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
	handle := goeyjs.CreateElement("textarea", "goey form-control")
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

func (w *textareaElement) createMeasurementElement() js.Value {
	handle := goeyjs.CreateElement("textarea", "goey-measure form-control")
	handle.Set("rows", w.minLines)

	goeyjs.AppendChildToBody(handle)

	return handle
}

func (w *textareaElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(bc.Max.Width)
	width = bc.ConstrainWidth(width)
	height := w.MinIntrinsicHeight(width)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *textareaElement) MinIntrinsicHeight(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	height := handle.Get("offsetHeight").Int()

	return base.FromPixelsY(height)
}

func (w *textareaElement) MinIntrinsicWidth(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := handle.Get("offsetWidth").Int()

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
