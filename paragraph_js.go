// +build go1.12

package goey

import (
	"fmt"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
)

type paragraphElement struct {
	Control
}

func (w *P) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "p")
	handle.Set("className", "goey")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &paragraphElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *paragraphElement) measureReflowLimits() {
	const textContent = "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm"

	handle := js.Global().Get("document").Call("createElement", "p")
	handle.Set("textContent", textContent)
	handle.Get("style").Set("visibility", "hidden")

	body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)
	width := handle.Get("offsetWidth").Int() + 1
	handle.Call("remove")

	paragraphMaxWidth = base.FromPixelsX(width)
}

func (w *paragraphElement) MinIntrinsicHeight(width base.Length) base.Length {
	if width == base.Inf {
		width = w.maxReflowWidth()
	}

	textContent := w.handle.Get("textContent")
	if textContent.String() == "" {
		textContent = js.ValueOf("X")
	}

	handle := js.Global().Get("document").Call("createElement", "p")
	handle.Set("textContent", textContent)
	handle.Get("style").Set("visibility", "hidden")
	handle.Get("style").Set("width", fmt.Sprintf("%dpx", width.PixelsX()))

	body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)
	height := handle.Get("offsetHeight").Int()
	handle.Call("remove")

	return base.FromPixelsY(height)
}

func (w *paragraphElement) MinIntrinsicWidth(height base.Length) base.Length {
	handle := js.Global().Get("document").Call("createElement", "p")
	handle.Set("textContent", w.handle.Get("innerText"))
	handle.Get("style").Set("visibility", "hidden")
	if height != base.Inf {
		handle.Get("style").Set("height", fmt.Sprintf("%dpx", height.PixelsY()))
	}

	body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)
	width := handle.Get("offsetWidth").Int()
	handle.Call("remove")

	return base.FromPixelsX(width)
}

func (w *paragraphElement) Props() base.Widget {
	getAlign := func(s string) TextAlignment {
		switch s {
		default:
			return JustifyLeft
		case "right":
			return JustifyRight
		case "center":
			return JustifyCenter
		case "justify":
			return JustifyFull
		}
	}

	return &P{
		Text:  w.handle.Get("textContent").String(),
		Align: getAlign(w.handle.Get("style").Get("text-align").String()),
	}
}

func (w *paragraphElement) updateProps(data *P) error {
	w.handle.Set("textContent", data.Text)
	switch data.Align {
	case JustifyLeft:
		w.handle.Get("style").Set("text-align", "left")
	case JustifyRight:
		w.handle.Get("style").Set("text-align", "right")
	case JustifyCenter:
		w.handle.Get("style").Set("text-align", "center")
	case JustifyFull:
		w.handle.Get("style").Set("text-align", "justify")
	}

	return nil
}
