//go:build go1.12
// +build go1.12

package goey

import (
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/js"
	"gitlab.com/stone.code/assert"
)

type labelElement struct {
	Control
}

func (w *Label) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := goeyjs.CreateElement("span", "goey")
	handle.Set("textContent", w.Text)
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &labelElement{
		Control: Control{handle},
	}

	return retval, nil
}

func (w *labelElement) createMeasurementElement() js.Value {
	text := w.handle.Get("textContent").String()
	if text == "" {
		text = "X"
	}

	handle := goeyjs.CreateElement("span", "goey-measure")
	handle.Set("textContent", text)

	goeyjs.AppendChildToBody(handle)

	return handle
}

func (w *labelElement) Layout(bc base.Constraints) base.Size {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := base.FromPixelsX(handle.Get("offsetWidth").Int() + 1)
	width = bc.ConstrainWidth(width)
	height := base.FromPixelsY(handle.Get("offsetHeight").Int() + 1)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *labelElement) MinIntrinsicHeight(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	height := handle.Get("offsetHeight").Int()
	assert.Assert(height > 0, "failure measuring label height")

	return base.FromPixelsY(height + 1)
}

func (w *labelElement) MinIntrinsicWidth(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := handle.Get("offsetWidth").Int()
	assert.Assert(width > 0, "failure measuring label width")

	return base.FromPixelsX(width + 1)
}

func (w *labelElement) Props() base.Widget {
	return &Label{
		Text: w.handle.Get("textContent").String(),
	}
}

func (w *labelElement) updateProps(data *Label) error {
	w.handle.Set("textContent", data.Text)

	return nil
}
