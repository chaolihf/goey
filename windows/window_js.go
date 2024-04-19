//go:build go1.12
// +build go1.12

package windows

import (
	"fmt"
	"image"
	"syscall/js"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/internal/js"
	"github.com/chaolihf/goey/loop"
	"gitlab.com/stone.code/assert"
)

type windowImpl struct {
	handle                  js.Value
	child                   base.Element
	horizontalScroll        bool
	horizontalScrollVisible bool
	verticalScroll          bool
	verticalScrollVisible   bool
	onClosing               func() bool
}

func init() {
	document := js.Global().Get("document")
	head := document.Call("getElementsByTagName", "head").Index(0)

	style := document.Call("createElement", "style")
	style.Set("textContent", `.goey {
		position:absolute; margin:0;
	}
	.goey-measure {
		position:absolute; visibility:hidden;
		width:auto; height:auto;
	}
	.goey-tabs-panel {
		border-left: solid 1px rgb(222,226,230);
		border-right: solid 1px rgb(222,226,230);
		border-bottom: solid 1px rgb(222,226,230);
	}`)

	head.Call("appendChild", style)
}

func newWindow(title string) (*Window, error) {
	document := js.Global().Get("document")
	document.Set("title", title)
	handle := document.Call("getElementsByTagName", "body").Index(0)
	assert.Assert(handle.Type() == js.TypeObject, "expected body element to be an object")

	// Clear body, if requested by the data attribute.
	if handle.Get("dataset").Get("goeyClear").Truthy() {
		println("clearing contents")
		handle.Set("innerHTML", "")
	}

	loop.AddLockCount(1)

	retval := &Window{windowImpl{
		handle: handle,
	}}

	js.Global().Get("window").Set("onresize", js.FuncOf(func(js.Value, []js.Value) interface{} {
		retval.onSize()
		return nil
	}))

	return retval, nil
}

func (w *windowImpl) OnDeleteEvent() bool {
	if w.onClosing == nil {
		return false
	}
	return w.onClosing()
}

func (w *windowImpl) control() base.Control {
	return base.Control{w.handle}
}

func (w *windowImpl) close() {
	if !w.handle.IsNull() {
		if w.child != nil {
			document := js.Global().Get("document")
			if ae := document.Get("activeElement"); ae.Truthy() {
				ae.Call("blur")
			}
			w.child.Close()
			w.child = nil
		}

		w.handle = js.Null()
		loop.AddLockCount(-1)
	}
}

func (w *windowImpl) onSize() {
	if w.child == nil {
		return
	}

	// Get the client area size.
	w.setDPI()
	clientSize := base.Size{
		base.FromPixelsX(js.Global().Get("window").Get("innerWidth").Int()),
		base.FromPixelsY(js.Global().Get("window").Get("innerHeight").Int()),
	}

	// Perform layout
	size := w.layoutChild(clientSize)
	bounds := base.Rectangle{
		base.Point{}, base.Point{size.Width, size.Height},
	}
	w.child.SetBounds(bounds)
}

func (w *windowImpl) setChildPost() {
	// Constrain window size
	w.updateWindowMinSize()
	// Properties may have changed sizes, so we need to do layout.
	w.onSize()
}

func (w *windowImpl) setScroll(horz, vert bool) {
	// If either scrollbar is being disabled, make sure that it is hidden.
	if !horz {
		w.handle.Get("style").Set("overflowX", "hidden")
		w.horizontalScrollVisible = false
	}
	if !vert {
		w.handle.Get("style").Set("overflowY", "hidden")
		w.verticalScrollVisible = false
	}

	// Redo layout to account for new box constraints, and show
	// scrollbars if necessary
	w.onSize()
}

func (w *windowImpl) show() {
}

func (w *windowImpl) showScrollH(width base.Length, clientWidth base.Length) bool {
	if width > clientWidth {
		if !w.horizontalScrollVisible {
			// Show the scrollbar
			w.handle.Get("style").Set("overflowX", "scroll")
			w.horizontalScrollVisible = true
			return true
		}
	} else if w.horizontalScrollVisible {
		// Remove the scroll bar
		w.handle.Get("style").Set("overflowX", "hidden")
		w.horizontalScrollVisible = false
		return true
	}

	return false
}

func (w *windowImpl) showScrollV(height base.Length, clientHeight base.Length) bool {
	if height > clientHeight {
		if !w.verticalScrollVisible {
			// Show the scrollbar
			w.handle.Get("style").Set("overflowY", "scroll")
			w.verticalScrollVisible = true
			return true
		}
	} else if w.verticalScrollVisible {
		// Remove the scroll bar
		w.handle.Get("style").Set("overflowY", "hidden")
		w.verticalScrollVisible = false
		return true
	}

	return false
}

func (_ *windowImpl) setDPI() {
	base.DPI.X, base.DPI.Y = 96, 96
}

func (w *windowImpl) setIcon(img image.Image) error {
	document := js.Global().Get("document")
	favicon := document.Call("getElementById", "goey-favicon")

	// If the link element does not yet exist, create it.
	if !favicon.Truthy() {
		favicon = document.Call("createElement", "link")
		favicon.Set("rel", "shortcut icon")
		favicon.Set("id", "goey-favicon")

		head := document.Call("getElementsByTagName", "head").Index(0)
		head.Call("appendChild", favicon)
	}

	// Set image data for the favicon.
	favicon.Set("href", goeyjs.ImageToAttr(img))

	return nil
}

func (w *windowImpl) setOnClosing(callback func() bool) {
	w.onClosing = callback
}

func (w *windowImpl) setTitle(value string) error {
	js.Global().Get("document").Set("title", value)
	return nil
}

func (w *windowImpl) title() string {
	return js.Global().Get("document").Get("title").String()
}

func (w *windowImpl) updateWindowMinSize() {
	size := w.MinSize()

	dx := size.Width.PixelsX()
	dy := size.Height.PixelsY()

	// Determine the extra width and height required for scrollbars.
	if w.verticalScroll {
		dx += 0 // int(gtk.WindowVScrollbarWidth(w.handle))
	}
	if w.horizontalScroll {
		dy += 0 // int(gtk.WindowHScrollbarHeight(w.handle))
	}

	style := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0).Get("style")
	style.Set("minWidth", fmt.Sprintf("%dpx", dx))
	style.Set("minHeight", fmt.Sprintf("%dpx", dy))
}
