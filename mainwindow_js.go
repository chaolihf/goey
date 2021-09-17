//go:build go1.12
// +build go1.12

package goey

import (
	"image"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/dialog"
	"bitbucket.org/rj/goey/loop"
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
	}`)

	head.Call("appendChild", style)
}

func newWindow(title string, child base.Widget) (*Window, error) {
	document := js.Global().Get("document")
	document.Set("title", title)
	handle := document.Call("getElementsByTagName", "body").Index(0)
	assert.Assert(handle.Type() == js.TypeObject, "expected body element to be an object")

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

func (w *windowImpl) message(m *dialog.Message) {
	// m.WithTitle(win2.GetWindowText(w.hWnd))
	// m.WithOwner(w.hWnd)
}

func (w *windowImpl) openfiledialog(m *dialog.OpenFile) {
	// m.WithTitle(win2.GetWindowText(w.hWnd))
	// m.WithOwner(w.hWnd)
}

func (w *windowImpl) savefiledialog(m *dialog.SaveFile) {
	// m.WithTitle(win2.GetWindowText(w.hWnd))
	// m.WithOwner(w.hWnd)
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
	// Redo the layout so the children are placed.
	if w.child != nil {
		// Constrain window size
		// w.updateWindowMinSize()
		// Properties may have changed sizes, so we need to do layout.
		w.onSize()
	} else {
		// Ensure that the scrollbars are hidden.
		//gtk.WindowShowScrollbars(w.handle, false, false)
	}
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
			//gtk.WindowShowScrollbars(w.handle, true, w.verticalScrollVisible)
			w.horizontalScrollVisible = true
			return true
		}
	} else if w.horizontalScrollVisible {
		// Remove the scroll bar
		//gtk.WindowShowScrollbars(w.handle, false, w.verticalScrollVisible)
		w.horizontalScrollVisible = false
		return true
	}

	return false
}

func (w *windowImpl) showScrollV(height base.Length, clientHeight base.Length) bool {
	if height > clientHeight {
		if !w.verticalScrollVisible {
			// Show the scrollbar
			//gtk.WindowShowScrollbars(w.handle, w.horizontalScrollVisible, true)
			w.verticalScrollVisible = true
			return true
		}
	} else if w.verticalScrollVisible {
		// Remove the scroll bar
		//gtk.WindowShowScrollbars(w.handle, w.horizontalScrollVisible, false)
		w.verticalScrollVisible = false
		return true
	}

	return false
}

func (_ *windowImpl) setDPI() {
	base.DPI.X, base.DPI.Y = 96, 96
}

func (w *windowImpl) setIcon(img image.Image) error {
	if img == nil {
		//gtk.WindowSetIcon(w.handle, nil, 0, 0, 0)
	}

	// rgba := imageToRGBA(img)
	//gtk.WindowSetIcon(w.handle, &rgba.Pix[0], rgba.Rect.Dx(), rgba.Rect.Dy(), rgba.Stride)
	//w.iconPix = rgba.Pix
	return nil
}

func (w *windowImpl) setOnClosing(callback func() bool) {
	w.onClosing = callback
}

func (w *windowImpl) setTitle(value string) error {
	js.Global().Get("document").Set("title", value)
	return nil
}

func (w *windowImpl) title() (string, error) {
	return js.Global().Get("document").Get("title").String(), nil
}

func (w *windowImpl) updateWindowMinSize() {
	// Determine the extra width and height required for borders, title bar,
	// and scrollbars
	dx, dy := 0, 0
	if w.verticalScroll {
		dx += 0 // int(gtk.WindowVScrollbarWidth(w.handle))
	}
	if w.horizontalScroll {
		dy += 0 // int(gtk.WindowHScrollbarHeight(w.handle))
	}

	// If there is no child, then we just need enough space for the window chrome.
	if w.child == nil {
		// gtk.WidgetSetSizeRequest(w.handle, dx, dy)
		return
	}

	request := image.Point{}
	// Determine the minimum size (in pixels) for the child of the window
	if w.horizontalScroll && w.verticalScroll {
		width := w.child.MinIntrinsicWidth(base.Inf)
		height := w.child.MinIntrinsicHeight(base.Inf)
		request.X = width.PixelsX() + dx
		request.Y = height.PixelsY() + dy
	} else if w.horizontalScroll {
		height := w.child.MinIntrinsicHeight(base.Inf)
		size := w.child.Layout(base.TightHeight(height))
		request.X = size.Width.PixelsX() + dx
		request.Y = height.PixelsY() + dy
	} else if w.verticalScroll {
		width := w.child.MinIntrinsicWidth(base.Inf)
		size := w.child.Layout(base.TightWidth(width))
		request.X = width.PixelsX() + dx
		request.Y = size.Height.PixelsY() + dy
	} else {
		width := w.child.MinIntrinsicWidth(base.Inf)
		height := w.child.MinIntrinsicHeight(base.Inf)
		size1 := w.child.Layout(base.TightWidth(width))
		size2 := w.child.Layout(base.TightHeight(height))
		request.X = max(width, size2.Width).PixelsX() + dx
		request.Y = max(height, size1.Height).PixelsY() + dy
	}

	// If scrolling is enabled for either direction, we can relax the
	// minimum window size.  These limits are fairly arbitrary, but we do need to
	// leave enough space for the scroll bars.
	if limit := (120 * DIP).PixelsX(); w.horizontalScroll && request.X > limit {
		request.X = limit
	}
	if limit := (120 * DIP).PixelsY(); w.verticalScroll && request.Y > limit {
		request.Y = limit
	}

	// gtk.WidgetSetSizeRequest(w.handle, request.X, request.Y)
}
