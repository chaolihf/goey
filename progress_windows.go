package goey

import (
	"syscall"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/lxn/win"
)

var (
	progress struct {
		className     *uint16
		oldWindowProc uintptr
	}
)

func init() {
	var err error
	progress.className, err = syscall.UTF16PtrFromString("msctls_progress32")
	if err != nil {
		panic(err)
	}
}

func (w *Progress) mount(parent base.Control) (base.Element, error) {
	style := uint32(win.WS_CHILD | win.WS_VISIBLE)
	hwnd := win.CreateWindowEx(0, progress.className, nil, style,
		10, 10, 100, 100,
		parent.HWnd, win.HMENU(nextControlID()), 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}
	win.SendMessage(hwnd, win.PBM_SETRANGE32, uintptr(w.Min), uintptr(w.Max))
	win.SendMessage(hwnd, win.PBM_SETPOS, uintptr(w.Value), 0)

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hMessageFont), 0)
	}

	retval := &progressElement{
		Control: Control{hwnd},
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type progressElement struct {
	Control
}

func (w *progressElement) Layout(bc base.Constraints) base.Size {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	width := w.MinIntrinsicWidth(0)
	if bc.Max.Width > 355*DIP {
		width = 355 * DIP
	}
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *progressElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 15 * DIP
}

func (w *progressElement) MinIntrinsicWidth(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 160 * DIP
}

func (w *progressElement) Props() base.Widget {
	min := win.SendMessage(w.hWnd, win.PBM_GETRANGE, win.TRUE, 0)
	max := win.SendMessage(w.hWnd, win.PBM_GETRANGE, win.FALSE, 0)
	value := win.SendMessage(w.hWnd, win.PBM_GETPOS, 0, 0)

	return &Progress{
		Value: int(value),
		Min:   int(min),
		Max:   int(max),
	}
}

func (w *progressElement) updateProps(data *Progress) error {
	win.SendMessage(w.hWnd, win.PBM_SETRANGE32, uintptr(data.Min), uintptr(data.Max))
	win.SendMessage(w.hWnd, win.PBM_SETPOS, uintptr(data.Value), 0)
	return nil
}
