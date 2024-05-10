package windows

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/chaolihf/goey/base"
)

var (
	// ErrSetChildrenNotReentrant is returned if a reentrant call to the method
	// SetChild is called.
	ErrSetChildrenNotReentrant = errors.New("method SetChild is not reentrant")

	insideSetChildren uintptr
)

// Window represents a top-level window that contain other widgets.
type Window struct {
	windowImpl
}

// NewWindow create a new top-level window for the application.
func NewWindow(title string, child base.Widget) (*Window, error) {
	// Create the window
	w, err := newWindow(title)
	if err != nil {
		return nil, err
	}

	// The the default values for the horizontal and vertical scroll.
	// We want to do this before creating the child so that scrollbars can
	// be displayed (if necessary) with the relayout for the child.
	w.horizontalScroll, w.verticalScroll = scrollDefaults()

	// Mount the widget, and initialize its layout.
	err = w.SetChild(child)
	if err != nil {
		w.Close()
		return nil, err
	}

	// Show the window
	w.show()

	if filename := os.Getenv("GOEY_SCREENSHOT"); filename != "" {
		asyncScreenshot(filename, w)
	}

	return w, nil
}

// Close destroys the window, and releases all associated resources.
func (w *Window) Close() {
	w.close()
}

// Child returns the mounted child for the window.  In general, this
// method should not be used.
func (w *Window) Child() base.Element {
	return w.child
}

func (w *windowImpl) layoutChild(windowSize base.Size) base.Size {
	// Create the constraints
	constraints := base.Tight(windowSize)

	// Relax maximum size when scrolling is allowed
	if w.horizontalScroll {
		constraints.Max.Width = base.Inf
	}
	if w.verticalScroll {
		constraints.Max.Height = base.Inf
	}

	// Perform layout
	size := w.child.Layout(constraints)
	if !constraints.IsSatisfiedBy(size) {
		fmt.Println("constraints not satisfied,", constraints, ",", size)
	}
	return size
}

// MinSize returns the minimum size required to layout the child.  The minimum
// size depends on the child, but also on what dimensions are allowed to scroll.
func (w *windowImpl) MinSize() base.Size {
	// In case the child needs to convert pixels to DIPs.
	w.setDPI()

	// Select strategy for calculating the minimum size depending on which
	// dimensions are allowed to scroll.
	if w.horizontalScroll && w.verticalScroll {
		return base.Size{
			Width:  min(w.child.MinIntrinsicWidth(base.Inf), 120*base.DIP),
			Height: min(w.child.MinIntrinsicHeight(base.Inf), 120*base.DIP),
		}
	} else if w.horizontalScroll {
		height := w.child.MinIntrinsicHeight(base.Inf)
		size := w.child.Layout(base.TightHeight(height))
		return base.Size{
			Width:  min(size.Width, 120*base.DIP),
			Height: height,
		}
	} else if w.verticalScroll {
		width := w.child.MinIntrinsicWidth(base.Inf)
		size := w.child.Layout(base.TightWidth(width))
		return base.Size{
			Width:  width,
			Height: min(size.Height, 120*base.DIP),
		}
	} else {
		width := w.child.MinIntrinsicWidth(base.Inf)
		height := w.child.MinIntrinsicHeight(base.Inf)
		size1 := w.child.Layout(base.TightWidth(width))
		size2 := w.child.Layout(base.TightHeight(height))
		return base.Size{
			Width:  max(width, size2.Width),
			Height: max(height, size1.Height),
		}
	}
}

// Scroll returns the flags that determine whether scrolling is allowed in the
// horizontal and vertical directions.
func (w *Window) Scroll() (horizontal, vertical bool) {
	return w.horizontalScroll, w.verticalScroll
}

func scrollDefaults() (horizontal, vertical bool) {
	env := os.Getenv("GOEY_SCROLL")
	if env == "" {
		return false, false
	}

	value, err := strconv.ParseUint(env, 10, 64)
	if err != nil || value >= 4 {
		return false, false
	}

	return (value & 2) == 2, (value & 1) == 1
}

// SetChild changes the child widget of the window.  As
// necessary, GUI widgets will be created or destroyed so that the GUI widgets
// match the widgets described by the parameter children.  The
// position of contained widgets will be updated to match the new layout
// properties.
func (w *Window) SetChild(child base.Widget) error {
	// One source of bugs in widgets is when the fire an event when being
	// updated.  This can lead to reentrant calls to SetChildren, typically
	// with incorrect information since the GUI is in an inconsistent state
	// when the event fires.  In short, this method is not reentrant.
	// The following will block changes to different windows, although
	// that shouldn't be susceptible to the same bugs.  Users in that
	// case should use Do to delay updates to other windows, but it shouldn't
	// happen in practice.
	if !atomic.CompareAndSwapUintptr(&insideSetChildren, 0, 1) {
		return ErrSetChildrenNotReentrant
	}
	defer func() {
		atomic.StoreUintptr(&insideSetChildren, 0)
	}()

	// The child may want to convert lengths to device dependent
	// pixels when mounting.
	w.setDPI()

	// Update the child element.
	newChild, err := base.DiffChild(w.control(), w.child, child)

	// Whether or not there has been an error, we need to run platform-specific
	// clean-up.  This is to recalculate min window size, update scrollbars, etc.
	w.child = newChild
	w.setChildPost()

	return err
}

// SetIcon changes the icon associated with the window.
//
// On Cocoa, individual windows do not have icons.  Instead, there is a single
// icon for the entire application.
func (w *Window) SetIcon(img image.Image) error {
	// Check that the image is not nil.  Want to enforce the precondition before
	// starting platform specific code to maintain uniformity.
	_ = img.Bounds()

	// Defer to platform specific code.
	return w.setIcon(img)
}

// SetOnClosing changes the event callback for when the user tries to close the
// window.  This callback can also be used to save or close any resources
// before the window is closed.
//
// Returning true from the callback will prevent the window from closing.
func (w *Window) SetOnClosing(callback func() bool) {
	w.setOnClosing(callback)
}

// SetScroll sets whether scrolling is allowed in the horizontal and vertical directions.
func (w *Window) SetScroll(horizontal, vertical bool) {
	// Copy the new parameters for the window into the fields.
	w.horizontalScroll, w.verticalScroll = horizontal, vertical

	// Defer to the platform dependent code.
	w.setScroll(horizontal, vertical)
}

// SetTitle changes the caption in the title bar for the window.
func (w *Window) SetTitle(title string) error {
	return w.setTitle(title)
}

// Title returns the current caption in the title bar for the window.
func (w *Window) Title() string {
	return w.title()
}

func sizeDefaults() (uint, uint) {
	const defaultWidth = 640
	const defaultHeight = 480

	env := os.Getenv("GOEY_SIZE")
	if env == "" {
		return defaultWidth, defaultHeight
	}

	parts := strings.Split(env, "x")
	if len(parts) != 2 {
		return defaultWidth, defaultHeight
	}

	width, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return defaultWidth, defaultHeight
	}

	height, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return defaultWidth, defaultHeight
	}

	return uint(width), uint(height)
}

// set callback for windows resize
func (w *Window) SetOnResize(callback func(int, int) bool) {
	w.setOnResize(callback)
}

// get windows client Size
func (w *Window) GetSize() (int, int) {
	return w.getSize()
}
