package goey

import (
	"github.com/lxn/win"
	"image"
	"syscall"
)

var (
	dpi image.Point
)

type DP int

type NativeWidget struct {
	hWnd win.HWND
}

func (w NativeWidget) Text() string {
	return GetWindowText(w.hWnd)
}

func (w NativeWidget) SetDisabled(value bool) {
	win.EnableWindow(w.hWnd, !value)
}

func (w *NativeWidget) SetBounds(bounds image.Rectangle) {

	win.MoveWindow(w.hWnd, int32(bounds.Min.X), int32(bounds.Min.Y), int32(bounds.Dx()), int32(bounds.Dy()), true)
}

func (w *NativeWidget) SetOrder(previous win.HWND) win.HWND {
	// Note, the argument previous may be 0 when setting the first child.
	// Fortunately, this corresponds to HWND_TOP, which sets the window
	// to top of the z-order.
	win.SetWindowPos(w.hWnd, previous, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
	return w.hWnd
}

func (w NativeWidget) SetText(value string) error {
	utf16, err := syscall.UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	rc := SetWindowText(w.hWnd, utf16)
	if rc == 0 {
		return syscall.GetLastError()
	}
	return nil
}

func (w *NativeWidget) Close() {
	if w.hWnd != 0 {
		win.DestroyWindow(w.hWnd)
		w.hWnd = 0
	}
}

type NativeMountedWidget interface {
	PreferredWidth() int
	CalculateHeight(width int) int
	SetBounds(bounds image.Rectangle)
	SetOrder(previous win.HWND) win.HWND
}
