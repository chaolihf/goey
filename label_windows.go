package goey

import (
	win2 "bitbucket.org/rj/goey/syscall"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	staticClassName     *uint16
	oldStaticWindowProc uintptr
)

func init() {
	var err error
	staticClassName, err = syscall.UTF16PtrFromString("STATIC")
	if err != nil {
		panic(err)
	}
}

func (w *Label) mount(parent Control) (Element, error) {
	text, err := syscall.UTF16FromString(w.Text)
	if err != nil {
		return nil, err
	}

	hwnd := win.CreateWindowEx(0, staticClassName, &text[0],
		win.WS_CHILD|win.WS_VISIBLE|win.SS_LEFT,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hMessageFont), 0)
	}

	retval := &mountedLabel{Control: Control{hwnd}, text: text}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedLabel struct {
	Control
	text []uint16
}

func (w *mountedLabel) preferredWidth() Length {
	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, 0x7fffffff, 0x7fffffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT, nil)
	win.ReleaseDC(w.hWnd, hdc)

	retval := FromPixelsX(int(rect.Right))
	return retval
}

func (w *mountedLabel) Props() Widget {
	return &Label{
		Text: w.Control.Text(),
	}
}

func (w *mountedLabel) Layout(bc Box) Size {
	// Determine ideal width.
	width := w.preferredWidth()
	height := 13 * DIP
	return bc.Constrain(Size{width, height})
}

func (w *mountedLabel) MinimumSize() Size {
	// Determine ideal width.
	width := w.preferredWidth()
	height := 13 * DIP
	return Size{width, height}
}

func (w *mountedLabel) SetBounds(bounds Rectangle) {
	w.Control.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *mountedLabel) updateProps(data *Label) error {
	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}
	w.text = text
	win2.SetWindowText(w.hWnd, &text[0])
	// TODO:  Update alignment

	return nil
}
