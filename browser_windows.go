package goey

import (
	"unsafe"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/win"
)

var (
	browser struct {
		className     []uint16
		oldWindowProc uintptr
	}
)

func init() {
	browser.className = []uint16{'B', 'R', 'O', 'W', 'S', 'E', 0}
}

func (w *Browser) mount(parent base.Control) (base.Element, error) {
	// Subclass the window procedure
	subclassWindowProcedure(w.EmbedWindow, &browser.oldWindowProc, browserWindowProc)

	retval := &browserElement{
		Control: Control{w.EmbedWindow},
	}
	win.SetWindowLongPtr(w.EmbedWindow, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type browserElement struct {
	Control
	text []uint16

	onClick func()
	onFocus func()
	onBlur  func()
}

func (w *browserElement) Props() base.Widget {
	return &Browser{
		EmbedWindow: w.Control.Hwnd,
	}
}

func (w *browserElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{Width: width, Height: height})
}

func (w *browserElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP
}

func (w *browserElement) MinIntrinsicWidth(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	width, _ := w.CalcRect(w.text)
	return max(
		75*DIP,
		base.FromPixelsX(int(width)+7),
	)
}

func (w *browserElement) updateProps(data *Browser) error {
	return nil
}

func browserWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		browserGetPtr(hwnd).Hwnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := browserGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := browserGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		// WM_COMMAND is sent to the parent, which will only forward certain
		// message.  This code should only ever see EN_UPDATE, but we will
		// still check.
		switch notification := win.HIWORD(uint32(wParam)); notification {
		case win.BN_CLICKED:
			if w := browserGetPtr(hwnd); w.onClick != nil {
				w.onClick()
			}
		}
		return 0
	case win.WM_PAINT:
		ps := win.PAINTSTRUCT{}
		win.BeginPaint(hwnd, &ps)
		win.EndPaint(hwnd, &ps)
		return 0
	}

	return win.CallWindowProc(browser.oldWindowProc, hwnd, msg, wParam, lParam)
}

func browserGetPtr(hwnd win.HWND) *browserElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*browserElement)(unsafe.Pointer(gwl))
	if ptr.Hwnd != hwnd && ptr.Hwnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
