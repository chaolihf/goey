package windows

import (
	"image"
	"syscall"
	"unsafe"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/dialog"
	win2 "github.com/chaolihf/goey/internal/windows"
	"github.com/chaolihf/goey/loop"
	"github.com/lxn/win"
)

var (
	className = []uint16{'G', 'o', 'e', 'y', 'M', 'a', 'i', 'n', 'W', 'i', 'n', 'd', 'o', 'w', 0}
	classAtom win.ATOM
)

type windowImpl struct {
	Hwnd                    win.HWND
	dpi                     image.Point
	windowRectDelta         image.Point
	windowMinSize           image.Point
	hicon                   win.HICON
	child                   base.Element
	childSize               base.Size
	onClosing               func() bool
	horizontalScroll        bool
	horizontalScrollVisible bool
	horizontalScrollPos     base.Length
	verticalScroll          bool
	verticalScrollVisible   bool
	verticalScrollPos       base.Length
}

func registerMainWindowClass() (win.ATOM, error) {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		return 0, syscall.GetLastError()
	}

	wc := win.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		HInstance:     hInstance,
		LpfnWndProc:   syscall.NewCallback(windowWindowProc),
		HCursor:       win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW)))),
		HbrBackground: win.GetSysColorBrush(win.COLOR_3DFACE),
		LpszClassName: &className[0],
	}

	atom := win.RegisterClassEx(&wc)
	if atom == 0 {
		return 0, syscall.GetLastError()
	}
	return atom, nil
}

func (w *windowImpl) onSize(hwnd win.HWND) {
	if w.child == nil {
		return
	}

	// Yes it's ugly, the SetBounds method for windows uses the screen DPI to
	// convert device independent pixels into actual pixels, but the DPI can change
	// from window to window when the computer has multiple monitors.  Fortunately,
	// all layout should happen in the GUI thread.
	w.setDPI()

	// Perform layout
	rect := win.RECT{}
	win.GetClientRect(hwnd, &rect)
	clientSize := base.FromPixels(
		int(rect.Right-rect.Left),
		int(rect.Bottom-rect.Top),
	)
	size := w.layoutChild(clientSize)

	// NOTE:  If the visibility of either scrollbar is changed, then a WM_SIZE
	// messagewill be sent, presumably because the size of the client area will
	// have changed.  This causes rentrant calls to onSize.  If a scrollbar is
	// either shown or hidden, then we need to abort layout.

	if w.horizontalScroll && w.verticalScroll {
		// Show scroll bars (both horizontal and vertical) if necessary.  Return
		// flag indicates whether visibility has changed.  We don't need to
		// worry about interaction of horizontal and vertical scrollbars, as any
		// change will force an abort and complete recalculation.
		ok := w.showScrollV(size.Height, clientSize.Height)
		if ok {
			return
		}
		ok = w.showScrollH(size.Width, clientSize.Width)
		if ok {
			return
		}
	} else if w.verticalScroll {
		// Show scroll bar if necessary.  Return flag indicates whether
		// visibility has changed.
		ok := w.showScrollV(size.Height, clientSize.Height)
		if ok {
			return
		}
	} else if w.horizontalScroll {
		// Show scroll bar if necessary.  Return flag indicates whether
		// visibility has changed.
		ok := w.showScrollH(size.Width, clientSize.Width)
		if ok {
			return
		}
	}
	w.childSize = size

	// Position the child element.
	w.child.SetBounds(base.Rectangle{
		Min: base.Point{-w.horizontalScrollPos, -w.verticalScrollPos},
		Max: base.Point{size.Width - w.horizontalScrollPos, size.Height - w.verticalScrollPos},
	})

	// Update the position of all of the children
	win.InvalidateRect(hwnd, &rect, true)
}

func newWindow(title string) (*Window, error) {
	//GetStartupInfo(&info);
	if win.OleInitialize() != win.S_OK {
		return nil, syscall.GetLastError()
	}

	// Ensure that our custom window class has been registered.
	if classAtom == 0 {
		atom, err := registerMainWindowClass()
		if err != nil {
			return nil, err
		}
		if atom == 0 {
			panic("internal error:  atom==0 although no error returned")
		}
		classAtom = atom
	}

	style := uint32(win.WS_OVERLAPPEDWINDOW | win.WS_CLIPCHILDREN)
	//if !settings.Resizable {
	//	style = win.WS_OVERLAPPED | win.WS_CAPTION | win.WS_MINIMIZEBOX | win.WS_SYSMENU
	//}

	rect := func() win.RECT {
		w, h := sizeDefaults()
		return win.RECT{
			Right:  int32(w),
			Bottom: int32(h),
		}
	}()
	win.AdjustWindowRect(&rect, win.WS_OVERLAPPEDWINDOW, false)

	var clientRect win.RECT
	win.GetClientRect(win2.GetDesktopWindow(), &clientRect)
	left := (clientRect.Right / 2) - ((rect.Right - rect.Left) / 2)
	top := (clientRect.Bottom / 2) - ((rect.Bottom - rect.Top) / 2)
	rect.Right = rect.Right - rect.Left + left
	rect.Left = left
	rect.Bottom = rect.Bottom - rect.Top + top
	rect.Top = top

	windowName, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return nil, err
	}
	hwnd := win.CreateWindowEx(win.WS_EX_CONTROLPARENT|win2.WS_EX_COMPOSITED, &className[0], windowName, style,
		rect.Left, rect.Top, rect.Right-rect.Left, rect.Bottom-rect.Top,
		win.HWND_DESKTOP, 0, 0, nil)
	if hwnd == 0 {
		win.OleUninitialize()
		return nil, syscall.GetLastError()
	}

	// Set the font for the window
	if hFont := win2.MessageFont(); hFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, 0, 0)
	}

	retval := &Window{windowImpl{Hwnd: hwnd}}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(&retval.windowImpl)))

	// Determine the DPI for this window
	hdc := win.GetDC(hwnd)
	retval.dpi.X = int(win.GetDeviceCaps(hdc, win.LOGPIXELSX))
	retval.dpi.Y = int(win.GetDeviceCaps(hdc, win.LOGPIXELSY))
	win.ReleaseDC(hwnd, hdc)

	// Calculate the extra width and height required for the borders
	windowRect := win.RECT{}
	win.GetWindowRect(hwnd, &windowRect)
	win.GetClientRect(hwnd, &clientRect)
	retval.windowRectDelta.X = int((windowRect.Right - windowRect.Left) - (clientRect.Right - clientRect.Left))
	retval.windowRectDelta.Y = int((windowRect.Bottom - windowRect.Top) - (clientRect.Bottom - clientRect.Top))

	return retval, nil
}

func (w *windowImpl) control() base.Control {
	return base.Control{w.Hwnd}
}

func (w *windowImpl) close() {
	// Want to be able to close windows in Go, even if they have already been
	// destroyed in the Win32 system
	if w.Hwnd != 0 {
		// There is a heseinbug with the kill focus message when destroying
		// windows.  To get consistent behavior, we can remove focus before
		// destroying the window.
		focus := win.GetFocus()
		if focus != 0 {
			parent := win.GetAncestor(focus, win.GA_ROOT)
			if parent == w.Hwnd {
				win.SetFocus(0)
			}
		}

		// Actually destroy the window.
		win.DestroyWindow(w.Hwnd)
		w.Hwnd = 0
	}

	// This call to uninitalize OLE is paired with a call in newWindow.
	win.OleUninitialize()
}

// NativeHandle returns the handle to the platform-specific window handle
// (i.e. a HWND on WIN32).
func (w *windowImpl) NativeHandle() win.HWND {
	return w.Hwnd
}

func (w *windowImpl) message(m *dialog.Message) {
	m.WithTitle(win2.GetWindowText(w.Hwnd))
	m.WithOwner(dialog.Owner{HWnd: w.Hwnd})
}

func (w *windowImpl) openfiledialog(m *dialog.OpenFile) {
	m.WithTitle(win2.GetWindowText(w.Hwnd))
	m.WithOwner(dialog.Owner{HWnd: w.Hwnd})
}

func (w *windowImpl) savefiledialog(m *dialog.SaveFile) {
	m.WithTitle(win2.GetWindowText(w.Hwnd))
	m.WithOwner(dialog.Owner{HWnd: w.Hwnd})
}

// Screenshot returns an image of the window, as displayed on screen.
func (w *windowImpl) Screenshot() (image.Image, error) {
	// Need the client rect for the window.
	region := win.RECT{}
	win.GetWindowRect(w.Hwnd, &region)

	// Create the device context and bitmap for the image
	hdcScreen := win.GetDC(0)
	defer func() {
		win.ReleaseDC(0, hdcScreen)
	}()
	hdc := win.CreateCompatibleDC(hdcScreen)
	defer func() {
		win.DeleteObject(win.HGDIOBJ(hdc))
	}()
	bitmap := win.CreateCompatibleBitmap(hdcScreen, region.Right-region.Left, region.Bottom-region.Top)
	defer func() {
		win.DeleteObject(win.HGDIOBJ(bitmap))
	}()
	win.SelectObject(hdc, win.HGDIOBJ(bitmap))
	rc := win.StretchBlt(hdc, 0, 0, region.Right-region.Left, region.Bottom-region.Top,
		hdcScreen, region.Left, region.Top, region.Right-region.Left, region.Bottom-region.Top,
		win.SRCCOPY)
	if !rc {
		err := syscall.GetLastError()
		if err == nil {
			err = syscall.EINVAL
		}
		return nil, err
	}

	// Convert the bitmap to a image.Image.
	img := win2.BitmapToImage(hdc, bitmap)
	return img, nil
}

// setChild updates the child element of the window.  It also updates any
// cached data linked to the child element, in particular the window's
// minimum size.  This function will also perform layout on the child.
func (w *windowImpl) setChildPost() {
	// Clear the cache of the minimum window size
	w.windowMinSize = image.Point{}

	// Ensure that tab-order is correct
	w.child.SetOrder(win.HWND_TOP)
	// Perform layout
	w.onSize(w.Hwnd)
}

func (w *windowImpl) setDPI() {
	base.DPI = image.Point{
		X: int(float32(w.dpi.X) * win2.MessageFontScale),
		Y: int(float32(w.dpi.Y) * win2.MessageFontScale),
	}
}

func (w *windowImpl) setScroll(hscroll, vscroll bool) {
	// If either scrollbar is being disabled, make sure that it is hidden.
	if !hscroll {
		w.horizontalScrollPos = 0
		w.horizontalScrollVisible = false
		win2.ShowScrollBar(w.Hwnd, win.SB_HORZ, win.FALSE)
	}
	if !vscroll {
		w.verticalScrollPos = 0
		w.verticalScrollVisible = false
		win2.ShowScrollBar(w.Hwnd, win.SB_VERT, win.FALSE)
	}

	// Changing the existence of scrollbar also changes the layout constraints.
	// Need to relayout the child.  If necessary, this will show the scrollbars.
	w.onSize(w.Hwnd)
}

func (w *windowImpl) setScrollPos(direction int32, wParam uintptr) {
	// Get all of the scroll bar information.
	si := win.SCROLLINFO{FMask: win.SIF_ALL}
	si.CbSize = uint32(unsafe.Sizeof(si))
	win.GetScrollInfo(w.Hwnd, direction, &si)

	// Save the position for comparison later on.
	currentPos := si.NPos
	switch win.LOWORD(uint32(wParam)) {
	// User clicked the HOME keyboard key.
	case win.SB_TOP:
		si.NPos = si.NMin

	// User clicked the END keyboard key.
	case win.SB_BOTTOM:
		si.NPos = si.NMax

	// User clicked the top or left arrow.
	case win.SB_LINEUP:
		if direction == win.SB_HORZ {
			si.NPos -= int32((13 * base.DIP).PixelsX())
		} else {
			si.NPos -= int32((13 * base.DIP).PixelsY())
		}

	// User clicked the bottom or right arrow.
	case win.SB_LINEDOWN:
		if direction == win.SB_HORZ {
			si.NPos += int32((13 * base.DIP).PixelsX())
		} else {
			si.NPos += int32((13 * base.DIP).PixelsY())
		}

	// User clicked the scroll bar shaft above or to the left of the scroll box.
	case win.SB_PAGEUP:
		si.NPos -= int32(si.NPage)

	// User clicked the scroll bar shaft below or to the right of the scroll
	// box.
	case win.SB_PAGEDOWN:
		si.NPos += int32(si.NPage)

	// User dragged the scroll box.
	case win.SB_THUMBTRACK:
		si.NPos = si.NTrackPos
	}

	// Set the position and then retrieve it.  Due to adjustments
	// by Windows it may not be the same as the value set.
	si.FMask = win.SIF_POS
	win.SetScrollInfo(w.Hwnd, direction, &si, true)
	win.GetScrollInfo(w.Hwnd, direction, &si)

	// If the position has changed, scroll window and update it.
	if si.NPos != currentPos {
		if direction == win.SB_HORZ {
			w.horizontalScrollPos = base.FromPixelsX(int(si.NPos))
		} else {
			w.verticalScrollPos = base.FromPixelsY(int(si.NPos))
		}
		w.child.SetBounds(base.Rectangle{
			Min: base.Point{-w.horizontalScrollPos, -w.verticalScrollPos},
			Max: base.Point{w.childSize.Width - w.horizontalScrollPos, w.childSize.Height - w.verticalScrollPos},
		})

		// TODO:  Use ScrollWindow function to reduce flicker during scrolling
		rect := win.RECT{}
		win.GetClientRect(w.Hwnd, &rect)
		win.InvalidateRect(w.Hwnd, &rect, true)
	}
}

func (w *windowImpl) show() {
	win.ShowWindow(w.Hwnd, win.SW_SHOW)
	win.UpdateWindow(w.Hwnd)
}

func (w *windowImpl) showScrollH(width base.Length, clientWidth base.Length) (flag bool) {
	if width > clientWidth {
		if !w.horizontalScrollVisible {
			// Create the scroll bar.  Any updates to the internal state must
			// be completed before the call, as this will send a WM_SIZE message
			// if the size of the client area changes.
			w.horizontalScrollVisible = true
			flag = true
			win2.ShowScrollBar(w.Hwnd, win.SB_HORZ, win.TRUE)
		}
		si := win.SCROLLINFO{
			FMask: win.SIF_PAGE | win.SIF_RANGE,
			NMin:  0,
			NMax:  int32(width.PixelsX()),
			NPage: uint32(clientWidth.PixelsX()),
		}
		si.CbSize = uint32(unsafe.Sizeof(si))
		win.SetScrollInfo(w.Hwnd, win.SB_HORZ, &si, true)
		si.FMask = win.SIF_POS
		win.GetScrollInfo(w.Hwnd, win.SB_HORZ, &si)
		w.horizontalScrollPos = base.FromPixelsX(int(si.NPos))
		return flag
	} else if w.horizontalScrollVisible {
		// Remove the scroll bar.
		w.horizontalScrollPos = 0
		w.horizontalScrollVisible = false
		win2.ShowScrollBar(w.Hwnd, win.SB_HORZ, win.FALSE)
		return true
	}

	return false
}

func (w *windowImpl) showScrollV(height base.Length, clientHeight base.Length) (flag bool) {
	if height > clientHeight {
		if !w.verticalScrollVisible {
			// Create the scroll bar.
			w.verticalScrollVisible = true
			flag = true
			win2.ShowScrollBar(w.Hwnd, win.SB_VERT, win.TRUE)
		}
		si := win.SCROLLINFO{
			FMask: win.SIF_PAGE | win.SIF_RANGE,
			NMin:  0,
			NMax:  int32(height.PixelsY()),
			NPage: uint32(clientHeight.PixelsY()),
		}
		si.CbSize = uint32(unsafe.Sizeof(si))
		win.SetScrollInfo(w.Hwnd, win.SB_VERT, &si, true)
		si.FMask = win.SIF_POS
		win.GetScrollInfo(w.Hwnd, win.SB_VERT, &si)
		w.verticalScrollPos = base.FromPixelsY(int(si.NPos))
		return flag
	} else if w.verticalScrollVisible {
		// Remove the scroll bar.
		w.verticalScrollPos = 0
		w.verticalScrollVisible = false
		win2.ShowScrollBar(w.Hwnd, win.SB_VERT, win.FALSE)
		return true
	}

	return false
}

func (w *windowImpl) setIcon(img image.Image) error {
	// Create the new icon
	hicon, err := win2.CreateIconFromImage(img)
	if err != nil {
		return err
	}

	// Delete the old hicon, and save the new hicon
	if w.hicon != 0 {
		win.DestroyIcon(w.hicon)
	}
	w.hicon = hicon

	// Update the window with the new icon.
	win2.SetClassLongPtr(w.Hwnd, win2.GCLP_HICON, uintptr(hicon))

	return nil
}

func (w *windowImpl) setOnClosing(callback func() bool) {
	w.onClosing = callback
}

func (w *windowImpl) setTitle(value string) error {
	_, err := win2.SetWindowText(w.Hwnd, value)
	return err
}

func (w *windowImpl) title() string {
	return win2.GetWindowText(w.Hwnd)
}

func (w *windowImpl) updateWindowMinSize() {
	// Determine the minimum client area size.
	size := w.MinSize()

	dx := size.Width.PixelsX()
	dy := size.Height.PixelsY()

	// Adjust the size to include space for scrollbars.
	if w.verticalScroll {
		// Want to include space for the scroll bar in the minimum width.
		// If the scrollbar is already visible, it will already be part
		// of the calculation through the difference in the window and client rectangles.
		dx += int(win.GetSystemMetrics(win.SM_CXVSCROLL))
	}
	if w.horizontalScroll {
		dy += int(win.GetSystemMetrics(win.SM_CYHSCROLL))
	}

	// Determine the extra width and height required for borders, title bar,
	// and scrollbars
	dx += w.windowRectDelta.X
	dy += w.windowRectDelta.Y

	w.windowMinSize.X = dx
	w.windowMinSize.Y = dy
}

func windowWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {

	switch msg {
	case win.WM_CREATE:
		// Maintain count of open windows.
		loop.AddLockCount(1)
		// Defer to default window proc

	case win.WM_NCDESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := windowGetPtr(hwnd); w != nil {
			w.Hwnd = 0
		}
		// Make sure we are no longer linked to as the active window
		loop.SetActiveWindow(0)
		// If this is the last main window visible, post the quit message so that the
		// message loop terminates.
		loop.AddLockCount(-1)
		// Defer to the default window proc

	case win.WM_CLOSE:
		if cb := windowGetPtr(hwnd).onClosing; cb != nil {
			if block := cb(); block {
				return 0
			}
		}
		// Defer to the default window proc

	case win.WM_ACTIVATE:
		if wParam != 0 {
			loop.SetActiveWindow(hwnd)
		}
		// Defer to the default window proc

	case win.WM_SETFOCUS:
		// The main window doesn't need focus, we want to delegate to a control
		if hwnd == win.GetFocus() { // Is this always true
			child := win.GetWindow(hwnd, win.GW_CHILD)
			for child != 0 {
				if style := win.GetWindowLong(child, win.GWL_STYLE); (style & win.WS_TABSTOP) != 0 {
					win.SetFocus(child)
					break
				}
				child = win.GetWindow(child, win.GW_HWNDNEXT)
			}
		}
		// Defer to the default window proc

	case win.WM_SIZE:
		windowGetPtr(hwnd).onSize(hwnd)
		// Defer to the default window proc

	case win.WM_GETMINMAXINFO:
		if w := windowGetPtr(hwnd); w != nil {
			if w.windowMinSize.X == 0 {
				w.updateWindowMinSize()
			}
			// Update tracking information based on our minimum size
			mmi := (*win.MINMAXINFO)(unsafe.Pointer(lParam))
			if limit := int32(w.windowMinSize.X); mmi.PtMinTrackSize.X < limit {
				mmi.PtMinTrackSize.X = limit
			}
			if limit := int32(w.windowMinSize.Y); mmi.PtMinTrackSize.Y < limit {
				mmi.PtMinTrackSize.Y = limit
			}
			return 0
		}
		// Defer to the default window proc

	case win.WM_HSCROLL:
		if lParam == 0 {
			// Message was sent by a standard scroll bar.  Need to adjust the
			// scroll position for the window.
			windowGetPtr(hwnd).setScrollPos(win.SB_HORZ, wParam)
		} else {
			// Message was sent by a child window.  As for all other controls
			// that notify the parent, resend to the child with the expectation
			// that the child has been subclassed.
			win.SendMessage(win.HWND(lParam), win.WM_HSCROLL, wParam, 0)
		}
		return 0

	case win.WM_VSCROLL:
		windowGetPtr(hwnd).setScrollPos(win.SB_VERT, wParam)
		return 0

	case win.WM_CTLCOLORSTATIC:
		win.SetBkMode(win.HDC(wParam), win.TRANSPARENT)
		return uintptr(win.GetSysColorBrush(win.COLOR_3DFACE))

	case win.WM_COMMAND:
		return WindowprocWmCommand(wParam, lParam)

	case win.WM_NOTIFY:
		n := (*win.NMHDR)(unsafe.Pointer(lParam))
		return win.SendMessage(n.HwndFrom, win.WM_NOTIFY, wParam, lParam)
	}

	// Let the default window proc handle all other messages
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func WindowprocWmCommand(wParam uintptr, lParam uintptr) uintptr {
	// These are the notifications that the controls needs to receive.
	if n := win.HIWORD(uint32(wParam)); n == win.BN_CLICKED || n == win.EN_UPDATE || n == win.CBN_SELCHANGE {
		// For BN_CLICKED, EN_UPDATE, and CBN_SELCHANGE, lParam is the window
		// handle of the control.  We don't need to use the control identifier
		// from wParam, we can dispatch directly to the control.
		return win.SendMessage(win.HWND(lParam), win.WM_COMMAND, wParam, lParam)
	}

	// Defer to the default window proc.  However, the default window proc will
	// return 0 for WM_COMMAND.
	return 0
}

func windowGetPtr(hwnd win.HWND) *windowImpl {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		return nil
	}

	ptr := (*windowImpl)(unsafe.Pointer(gwl))
	if ptr.Hwnd != hwnd && ptr.Hwnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
