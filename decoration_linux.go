package goey

import (
	"image/color"
	"math"
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

func (w *Decoration) mount(parent Control) (Element, error) {
	control, err := gtk.DrawingAreaNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &decorationElement{
		handle: control,
		fill:   w.Fill,
		insets: w.Insets,
		radius: w.Radius,
	}

	control.Connect("destroy", decoration_onDestroy, retval)
	control.Connect("draw", decoration_onDraw, retval)
	control.Show()

	child, err := DiffChild(parent, nil, w.Child)
	if err != nil {
		control.Destroy()
		return nil, err
	}
	retval.child = child

	return retval, nil
}

type decorationElement struct {
	handle *gtk.DrawingArea
	fill   color.RGBA
	insets Insets
	radius Length

	child     Element
	childSize Size
}

func decoration_onDestroy(widget *gtk.DrawingArea, mounted *decorationElement) {
	mounted.handle = nil
}

func decoration_onDraw(widget *gtk.DrawingArea, cr *cairo.Context, mounted *decorationElement) bool {
	a := mounted.handle.GetAllocation()
	if mounted.radius > 0 {
		radius := float64(mounted.radius.PixelsX())
		w, h := float64(a.GetWidth()), float64(a.GetHeight())
		if 2*radius > w {
			radius = w / 2
		}
		if 2*radius > h {
			radius = h / 2
		}
		cr.MoveTo(0, radius)
		cr.Arc(radius, radius, radius, math.Pi, 3*math.Pi/2)
		cr.LineTo(w-radius, 0)
		cr.Arc(w-radius, radius, radius, 3*math.Pi/2, 2*math.Pi)
		cr.LineTo(w, h-radius)
		cr.Arc(w-radius, h-radius, radius, 0, math.Pi/2)
		cr.LineTo(radius, h)
		cr.Arc(radius, h-radius, radius, math.Pi/2, math.Pi)
		cr.ClosePath()
	} else {
		cr.Rectangle(0, 0, float64(a.GetWidth()), float64(a.GetHeight()))
	}
	cr.SetSourceRGB(float64(mounted.fill.R)/0xFF, float64(mounted.fill.G)/0xFF, float64(mounted.fill.B)/0xFF)
	cr.Fill()
	return false
}

func (w *decorationElement) Close() {
	if w.handle != nil {
		w.child.Close()
		w.child = nil
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *decorationElement) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())

	bounds.Min.X += w.insets.Left
	bounds.Min.Y += w.insets.Top
	bounds.Max.X -= w.insets.Right
	bounds.Max.Y -= w.insets.Bottom
	w.child.SetBounds(bounds)
}

func (w *decorationElement) updateProps(data *Decoration) error {
	w.fill = data.Fill
	w.radius = data.Radius

	parent, err := w.handle.GetParent()
	if err != nil {
		return err
	}
	w.child, err = DiffChild(Control{parent}, w.child, data.Child)
	if err != nil {
		return err
	}

	return nil
}
