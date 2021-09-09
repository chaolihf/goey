package goey

import (
	"fmt"
	"image/color"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"gitlab.com/stone.code/assert"
)

type decorationElement struct {
	Control

	child     base.Element
	childSize base.Size
	insets    Insets
}

func (w *Decoration) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "div")
	handle.Set("className", "panel panel-default")
	handle.Get("style").Set("position", "absolute")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &decorationElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *decorationElement) Close() {
	w.child.Close()
	w.Control.Close()
}

func (w *decorationElement) props() *Decoration {
	println(w.handle.Get("style").Get("borderRadius").String())
	return &Decoration{}
}

func (w *decorationElement) SetBounds(bounds base.Rectangle) {
	// Update background control position
	w.Control.SetBounds(bounds)

	px := base.FromPixelsX(1)
	py := base.FromPixelsY(1)
	position := bounds.Min
	bounds.Min.X += px + w.insets.Left - position.X
	bounds.Min.Y += py + w.insets.Top - position.Y
	bounds.Max.X -= px + w.insets.Right + position.X
	bounds.Max.Y -= py + w.insets.Bottom + position.Y
	w.child.SetBounds(bounds)
}

func (w *decorationElement) updateProps(data *Decoration) (err error) {
	w.child, err = base.DiffChild(base.Control{w.handle}, w.child, data.Child)
	if err != nil {
		return err
	}
	assert.Assert(w.child != nil, "child is nil, but no error mountingn child")

	w.insets = data.Insets

	w.handle.Get("style").Set("borderRadius", fmt.Sprintf("%dpx", data.Radius.PixelsX()))
	if data.Stroke.A != 0 {
		w.handle.Get("style").Set("borderColor", cssColor(data.Stroke))
		w.handle.Get("style").Set("borderStyle", "solid")
		w.handle.Get("style").Set("borderWidth", "1px")
	}
	if data.Fill.A != 0 {
		w.handle.Get("style").Set("backgroundColor", cssColor(data.Fill))
	}

	return nil
}

func cssColor(clr color.RGBA) string {
	if clr.A == 0xFF {
		return fmt.Sprintf("#%02x%02x%02x", clr.R, clr.G, clr.B)
	}

	return fmt.Sprintf("#%02x%02x%02x%02x", clr.R, clr.G, clr.B, clr.A)
}
