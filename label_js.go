package goey

import (
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"gitlab.com/stone.code/assert"
)

type labelElement struct {
	Control
}

func (w *Label) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "span")
	handle.Set("innerText", w.Text)
	handle.Get("style").Set("position", "absolute")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &labelElement{
		Control: Control{handle},
	}

	return retval, nil
}

func (w *labelElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(bc.Max.Width)
	width = bc.ConstrainWidth(width)
	height := w.MinIntrinsicHeight(width)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *labelElement) MinIntrinsicHeight(base.Length) base.Length {
	text := w.handle.Get("innerText").String()
	if text == "" {
		text = "X"
	}

	// Create a dummy button
	handle := js.Global().Get("document").Call("createElement", "span")
	handle.Set("innerText", text)
	handle.Get("style").Set("visibility", "hidden")
	handle.Get("style").Set("display", "block")

	body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)
	height := handle.Get("offsetHeight").Int()
	assert.Assert(height > 0, "failure measuring label height")
	handle.Call("remove")

	return base.FromPixelsY(height)
}

func (w *labelElement) MinIntrinsicWidth(base.Length) base.Length {
	// Create a dummy button
	handle := js.Global().Get("document").Call("createElement", "span")
	handle.Set("innerText", w.handle.Get("innerText"))
	handle.Get("style").Set("visibility", "hidden")
	handle.Get("style").Set("display", "block")

	body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)
	width := handle.Get("offsetWidth").Int()
	assert.Assert(width > 0, "failure measuring label width")
	handle.Call("remove")

	return base.FromPixelsX(width)
}

func (w *labelElement) Props() base.Widget {
	return &Label{
		Text: w.handle.Get("innerText").String(),
	}
}

func (w *labelElement) updateProps(data *Label) error {
	w.handle.Set("innerText", data.Text)

	return nil
}
