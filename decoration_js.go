// +build go1.12

package goey

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
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
	if err := retval.updateProps(w); err != nil {
		retval.Control.Close()
		return nil, err
	}

	assert.Assert(retval.child != nil, "child should not be nil")
	return retval, nil
}

func (w *decorationElement) Close() {
	w.child.Close()
	w.Control.Close()
}

func (w *decorationElement) props() *Decoration {
	fromCSS := func(s string) color.RGBA {
		if s == "inherit" || s == "none" || s == "" {
			return color.RGBA{}
		}
		if ndx := strings.Index(s, "rgb"); ndx >= 0 {
			s = s[ndx:]

			if strings.Count(s, ",") == 2 {
				var clr color.RGBA
				fmt.Sscanf(s, "rgb(%d,%d,%d)", &clr.R, &clr.G, &clr.B)
				clr.A = 0xff
				return clr
			} else {
				var clr color.RGBA
				fmt.Sscanf(s, "rgb(%d,%d,%d,%d)", &clr.R, &clr.G, &clr.B, &clr.A)
				return clr
			}
		}

		panic("not implemented: " + s)
	}
	radiusFromCSS := func(s string) base.Length {
		if strings.HasSuffix(s, "px") {
			v, _ := strconv.Atoi(s[:len(s)-2])
			return base.FromPixelsX(v)
		}

		panic("not implemented: " + s)
	}

	return &Decoration{
		Insets: w.insets,
		Fill:   fromCSS(w.handle.Get("style").Get("background").String()),
		Stroke: fromCSS(w.handle.Get("style").Get("border").String()),
		Radius: radiusFromCSS(w.handle.Get("style").Get("borderRadius").String()),
	}
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

	style := w.handle.Get("style")
	style.Set("borderRadius", fmt.Sprintf("%dpx", data.Radius.PixelsX()))
	if data.Stroke.A != 0 {
		style.Set("border", cssColor(data.Stroke)+" solid 1px")
	} else {
		style.Set("border", "none")
	}
	if data.Fill.A != 0 {
		style.Set("background", cssColor(data.Fill))
	} else {
		style.Set("background", "inherit")
	}

	return nil
}

func cssColor(clr color.RGBA) string {
	if clr.A == 0xFF {
		return fmt.Sprintf("#%02x%02x%02x", clr.R, clr.G, clr.B)
	}

	return fmt.Sprintf("#%02x%02x%02x%02x", clr.R, clr.G, clr.B, clr.A)
}
