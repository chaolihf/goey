package goey

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/chaolihf/goey/base"
	win2 "github.com/chaolihf/goey/internal/windows"
	"github.com/chaolihf/goey/loop"
	"github.com/chaolihf/win"
)

func init() {
	// If the return of the call to InitCommonControlsEx is checked, we see
	// false, which according to the documentation indicates that it failed.
	// However, there is no error with syscall.GetLastError().
	//
	// Note:  The init function for github.com/chaolihf/win also calls this
	// function, but does not include ICC_STANDARD_CLASSES.
	initCtrls := win.INITCOMMONCONTROLSEX{}
	initCtrls.DwSize = uint32(unsafe.Sizeof(initCtrls))
	initCtrls.DwICC = win.ICC_STANDARD_CLASSES | win.ICC_DATE_CLASSES | win.ICC_TAB_CLASSES
	win.InitCommonControlsEx(&initCtrls)
}

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkWidget when targeting GTK.
//
// Unless developping new widgets, users should not need to use this type.
//
// Any method's on this type will be platform specific.
type Control struct {
	Hwnd win.HWND
}

// Text copies text of the underlying window.
func (w Control) Text() string {
	return win2.GetWindowText(w.Hwnd)
}

// CalcRect is a wrapper around the WIN32 call DrawTextEx with the option DT_CALCRECT.
func (w Control) CalcRect(text []uint16) (int32, int32) {
	hdc := win.GetDC(w.Hwnd)
	if hFont := win2.MessageFont(); hFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hFont))
	}

	rect := win.RECT{0, 0, 0x7fffffff, 0x7fffffff}
	win.DrawTextEx(hdc, &text[0], int32(len(text)), &rect, win.DT_CALCRECT, nil)
	win.ReleaseDC(w.Hwnd, hdc)

	return rect.Right, rect.Bottom
}

// SetDisabled is a wrapper around the WIN32 call to EnableWindow.
func (w Control) SetDisabled(value bool) {
	win.EnableWindow(w.Hwnd, !value)
}

// SetBounds is a wrapper around the WIN32 call to MoveWindow.
func (w *Control) SetBounds(bounds base.Rectangle) {
	win.MoveWindow(w.Hwnd, int32(bounds.Min.X.PixelsX()), int32(bounds.Min.Y.PixelsY()), int32(bounds.Dx().PixelsX()), int32(bounds.Dy().PixelsY()), false)
}

// TakeFocus is a wrapper around SetFocus.
func (w *Control) TakeFocus() bool {
	// If the control already has focus, we avoid the call to SetFocus.  This
	// is to debounce the event callbacks.
	if win.GetFocus() == w.Hwnd {
		return true
	}

	return win.SetFocus(w.Hwnd) != 0
}

// TypeKeys sends events to the control as if the string was typed by a user.
func (w *Control) TypeKeys(text string) chan error {
	errc := make(chan error, 1)

	go func() {
		defer close(errc)

		time.Sleep(50 * time.Millisecond)

		err := loop.Do(func() error {
			if win.GetForegroundWindow() == 0 {
				return fmt.Errorf("can't type keys: no foreground window")
			}
			return nil
		})
		if err != nil {
			errc <- err
			return
		}

		for _, r := range text {
			inp := [2]win.KEYBD_INPUT{
				{Type: win.INPUT_KEYBOARD, Ki: win.KEYBDINPUT{}},
				{Type: win.INPUT_KEYBOARD, Ki: win.KEYBDINPUT{}},
			}

			if r == '\n' {
				inp[0].Ki.WVk = win.VK_RETURN
				inp[1].Ki.WVk = win.VK_RETURN
				inp[1].Ki.DwFlags = win.KEYEVENTF_KEYUP
			} else {
				inp[0].Ki.WScan = uint16(r)
				inp[0].Ki.DwFlags = win.KEYEVENTF_UNICODE
				inp[1].Ki.WScan = uint16(r)
				inp[1].Ki.DwFlags = win.KEYEVENTF_UNICODE | win.KEYEVENTF_KEYUP
			}

			err := loop.Do(func() error {
				rc := win.SendInput(2, unsafe.Pointer(&inp), int32(unsafe.Sizeof(inp[0])))
				if rc != 2 {
					return fmt.Errorf("failed to send input: rc= %d: %x", rc, win.GetLastError())
				}
				return nil
			})
			if err != nil {
				errc <- err
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	return errc
}

// SetOrder is a call around SetWindowPos used to ensure that a window appears
// in the correct order.
func (w *Control) SetOrder(previous win.HWND) win.HWND {
	// Note, the argument previous may be 0 when setting the first child.
	// Fortunately, this corresponds to HWND_TOP, which sets the window
	// to top of the z-order.
	win.SetWindowPos(w.Hwnd, previous, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_NOREDRAW|0x400)
	return w.Hwnd
}

// Close is a wrapper around the WIN32 call to DestroyWindow.
func (w *Control) Close() {
	if w.Hwnd != 0 {
		win.DestroyWindow(w.Hwnd)
		w.Hwnd = 0
	}
}

func createControlWindow(exStyle uint32, classname *uint16, text string, style uint32, parent win.HWND) (win.HWND, []uint16, error) {
	// Get the text for the control.  There may be extra work here if the
	// string is empty, but that is not expected to be common.
	utftext, err := syscall.UTF16FromString(text)
	if err != nil {
		return 0, nil, err
	}

	// Create the control.
	hwnd := win.CreateWindowEx(exStyle, classname, &utftext[0], style,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		parent, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return 0, nil, syscall.EINVAL
		}
		return 0, nil, err
	}

	// Set the font for the window
	if hFont := win2.MessageFont(); hFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hFont), 0)
	}

	return hwnd, utftext, nil
}

func subclassWindowProcedure(hWnd win.HWND, oldWindowProc *uintptr, newWindowProc func(win.HWND, uint32, uintptr, uintptr) uintptr) {
	// We need a copy of the address of the old window proc when subclassing.
	// Unhandled messages need to be forwarded.
	if *oldWindowProc == 0 {
		*oldWindowProc = win.GetWindowLongPtr(hWnd, win.GWLP_WNDPROC)
	} else {
		// Paranoia.  Windows created with the same class should have the same
		// window proc set, but just in case we will double check.
		tmp := win.GetWindowLongPtr(hWnd, win.GWLP_WNDPROC)
		if tmp != *oldWindowProc {
			panic("Window procedure does not match.")
		}
	}

	// Subclass the window by setting a new window proc.
	win.SetWindowLongPtr(hWnd, win.GWLP_WNDPROC, syscall.NewCallback(newWindowProc))
}
