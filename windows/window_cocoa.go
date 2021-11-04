//go:build cocoa || (darwin && !gtk)
// +build cocoa darwin,!gtk

package windows

import (
	"image"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/dialog"
	"bitbucket.org/rj/goey/internal/cocoa"
	"bitbucket.org/rj/goey/loop"
)

type windowImpl struct {
	handle                  *cocoa.Window
	contentView             *cocoa.View
	child                   base.Element
	horizontalScroll        bool
	horizontalScrollVisible bool
	verticalScroll          bool
	verticalScrollVisible   bool

	onClosing func() bool
}

func newWindow(title string) (*Window, error) {
	// Update the global DPI
	base.DPI.X, base.DPI.Y = 96, 96

	w, h := sizeDefaults()
	handle := cocoa.NewWindow(title, w, h)
	loop.AddLockCount(1)
	retval := &Window{windowImpl{
		handle:      handle,
		contentView: handle.ContentView(),
	}}
	handle.SetCallbacks((*windowCallbacks)(&retval.windowImpl))

	return retval, nil
}

func (w *windowImpl) control() base.Control {
	return base.Control{w.contentView}
}

func (w *windowImpl) close() {
	if w.handle != nil {
		w.handle.Close()
		w.handle = nil
	}
}

func (w *windowImpl) message(m *dialog.Message) {
	//m.title, m.err = w.handle.GetTitle()
	m.WithOwner(dialog.Owner{Window: w.handle})
}

func (w *windowImpl) openfiledialog(m *dialog.OpenFile) {
	//m.title, m.err = w.handle.GetTitle()
	m.WithOwner(dialog.Owner{Window: w.handle})
}

func (w *windowImpl) savefiledialog(m *dialog.SaveFile) {
	//m.title, m.err = w.handle.GetTitle()
	m.WithOwner(dialog.Owner{Window: w.handle})
}

func (w *windowImpl) onSize() {
	if w.child == nil {
		return
	}

	// Update the global DPI
	w.setDPI()

	// Calculate the layout.
	width, height := w.handle.ContentSize()
	clientSize := base.Size{base.FromPixelsX(width), base.FromPixelsY(height)}
	size := w.layoutChild(clientSize)
	if w.horizontalScroll && w.verticalScroll {
		// Show scroll bars if necessary.
		w.showScrollV(size.Height, clientSize.Height)
		ok := w.showScrollH(size.Width, clientSize.Width)
		// Adding horizontal scroll take vertical space, so we need to check
		// again for vertical scroll.
		if ok {
			_, height := w.handle.ContentSize()
			w.showScrollV(size.Height, base.FromPixelsY(height))
		}
	} else if w.verticalScroll {
		// Show scroll bars if necessary.
		ok := w.showScrollV(size.Height, clientSize.Height)
		if ok {
			width, height := w.handle.ContentSize()
			clientSize := base.Size{base.FromPixelsX(width), base.FromPixelsY(height)}
			size = w.layoutChild(clientSize)
		}
	} else if w.horizontalScroll {
		// Show scroll bars if necessary.
		ok := w.showScrollH(size.Width, clientSize.Width)
		if ok {
			width, height := w.handle.ContentSize()
			clientSize := base.Size{base.FromPixelsX(width), base.FromPixelsY(height)}
			size = w.layoutChild(clientSize)
		}
	}
	w.handle.SetContentSize(int(size.Width.PixelsX()), int(size.Height.PixelsY()))

	// Set bounds on child control.
	bounds := base.Rectangle{
		base.Point{}, base.Point{size.Width, size.Height},
	}
	w.child.SetBounds(bounds)
}

// Screenshot returns an image of the window, as displayed on screen.
func (w *windowImpl) Screenshot() (image.Image, error) {
	img := w.handle.Screenshot()
	return img, nil
}

func (w *windowImpl) setChildPost() {
	// Constrain window size
	w.updateWindowMinSize()
	// Properties may have changed sizes, so we need to do layout.
	w.onSize()
}

// setDPI updates the global DPI.
func (*windowImpl) setDPI() {
	base.DPI.X, base.DPI.Y = 96, 96
}

func (w *windowImpl) setScroll(horz, vert bool) {
	// If either scrollbar is being disabled, make that it is hidden.
	if !horz || !vert {
		w.handle.SetScrollVisible(false, false)
		w.horizontalScrollVisible = false
		w.verticalScrollVisible = false
	}

	// Redo layout to account for new box constraints, and show
	// scrollbars if necessary
	w.onSize()
}

func (w *windowImpl) show() {
	//w.handle.ShowAll()
}

func (w *windowImpl) setIcon(img image.Image) error {
	return w.handle.SetIcon(img)
}

func (w *windowImpl) setOnClosing(callback func() bool) {
	w.onClosing = callback
}

func (w *windowImpl) setTitle(value string) error {
	w.handle.SetTitle(value)
	return nil
}

func (w *windowImpl) showScrollH(width base.Length, clientWidth base.Length) bool {
	if width > clientWidth {
		if !w.horizontalScrollVisible {
			// Show the scrollbar
			w.handle.SetScrollVisible(true, w.verticalScrollVisible)
			w.horizontalScrollVisible = true
			return true
		}
	} else if w.horizontalScrollVisible {
		// Remove the scroll bar
		w.handle.SetScrollVisible(false, w.verticalScrollVisible)
		w.horizontalScrollVisible = false
		return true
	}

	return false
}

func (w *windowImpl) showScrollV(height base.Length, clientHeight base.Length) bool {
	if height > clientHeight {
		if !w.verticalScrollVisible {
			// Show the scrollbar
			w.handle.SetScrollVisible(w.horizontalScrollVisible, true)
			w.verticalScrollVisible = true
			return true
		}
	} else if w.verticalScrollVisible {
		// Remove the scroll bar
		w.handle.SetScrollVisible(w.horizontalScrollVisible, false)
		w.verticalScrollVisible = false
		return true
	}

	return false
}

func (w *windowImpl) title() string {
	return w.handle.Title()
}

func (w *windowImpl) updateWindowMinSize() {
	size := w.MinSize()

	dx := size.Width.PixelsX()
	dy := size.Height.PixelsY()

	// Determine the extra width and height required for scrollbars.
	if w.verticalScroll {
		// TODO:  Measure scrollbar width
		dx += 15
	}
	if w.horizontalScroll {
		// TODO:  Measure scrollbar height
		dy += 15
	}

	w.handle.SetMinSize(dx, dy)
}

type windowCallbacks windowImpl

func (w *windowCallbacks) OnShouldClose() bool {
	if w.onClosing != nil {
		return !w.onClosing()
	}
	return true
}

func (w *windowCallbacks) OnWillClose() {
	w.handle = nil
	loop.AddLockCount(-1)
}

func (w *windowCallbacks) OnDidResize() {
	impl := (*windowImpl)(w)
	impl.onSize()
}
