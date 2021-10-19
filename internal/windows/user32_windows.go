package windows

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var (
	moduser32 = syscall.MustLoadDLL("user32.dll")

	procSetClassLongPtr     = moduser32.MustFindProc("SetClassLongPtrW")
	procGetDesktopWindow    = moduser32.MustFindProc("GetDesktopWindow")
	procGetWindowText       = moduser32.MustFindProc("GetWindowTextW")
	procGetWindowTextLength = moduser32.MustFindProc("GetWindowTextLengthW")
	procSetWindowText       = moduser32.MustFindProc("SetWindowTextW")
	procShowScrollBar       = moduser32.MustFindProc("ShowScrollBar")
)

const (
	GCLP_HICON   = -14
	GCLP_HICONSM = -34

	DTM_FIRST         = 0x1000
	DTM_CLOSEMONTHCAL = DTM_FIRST + 13

	MCM_FIRST  = 0x1000
	MCN_FIRST  = uint32(0xFFFFFD12)
	MCN_SELECT = MCN_FIRST + 4

	STM_SETIMAGE = 0x0172
	STM_GETIMAGE = 0x0173

	TBS_HORZ        = 0x0000
	TBS_AUTOTICKS   = 0x0001
	TBM_SETTICFREQ  = win.WM_USER + 20
	TBM_SETPAGESIZE = win.WM_USER + 21
	TBM_SETLINESIZE = win.WM_USER + 23

	WS_EX_COMPOSITED = 0x02000000
)

// NMSELCHANGE match the C structure of the same name.
type NMSELCHANGE struct {
	Nmhdr      win.NMHDR
	StSelStart win.SYSTEMTIME
	StSelEnd   win.SYSTEMTIME
}

// SetClassLongPtr is a wrapper.
func SetClassLongPtr(hWnd win.HWND, index int32, value uintptr) uintptr {
	ret, _, _ := syscall.Syscall(procSetClassLongPtr.Addr(), 3,
		uintptr(hWnd),
		uintptr(index),
		value)

	return ret
}

// GetDesktopWindow is a wrapper.
func GetDesktopWindow() win.HWND {
	r1, _, err := syscall.Syscall(procGetDesktopWindow.Addr(), 0, 0, 0, 0)
	if err != 0 {
		panic(err)
	}
	return win.HWND(r1)
}

// GetWindowText is a wrapper for GetWindowTextLength and GetWindowText.
// This function provides a somewhat higher-level API than the C API, as Go
// is garbage collected, so the buffer management provided by the C API is
// not required.
func GetWindowText(hWnd win.HWND) string {
	r0, _, _ := syscall.Syscall(procGetWindowTextLength.Addr(), 1, uintptr(hWnd), 0, 0)
	if r0 < 80 {
		var buffer [80]uint16
		r0, _, _ := syscall.Syscall(procGetWindowText.Addr(), 3, uintptr(hWnd), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
		return syscall.UTF16ToString(buffer[:r0])
	}
	buffer := make([]uint16, r0)
	r0, _, _ = syscall.Syscall(procGetWindowText.Addr(), 3, uintptr(hWnd), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	return syscall.UTF16ToString(buffer[:r0])
}

// GetWindowTextLength is a wrapper.
func GetWindowTextLength(hWnd win.HWND) int32 {
	r0, _, _ := syscall.Syscall(procGetWindowTextLength.Addr(), 1, uintptr(hWnd), 0, 0)
	return int32(r0)
}

// SetWindowText is a wrapper.
func SetWindowText(hWnd win.HWND, text string) ([]uint16, error) {
	utf16, err := syscall.UTF16FromString(text)
	if err != nil {
		return nil, err
	}

	_, _, errno := syscall.Syscall(procSetWindowText.Addr(), 2, uintptr(hWnd), uintptr(unsafe.Pointer(&utf16[0])), 0)
	if errno != 0 {
		return utf16, errno
	}
	return utf16, nil
}

// ShowScrollBar is a wrapper.
func ShowScrollBar(hWnd win.HWND, wSBFlags uint, bShow win.BOOL) win.BOOL {
	r0, _, _ := syscall.Syscall(procShowScrollBar.Addr(), 3, uintptr(hWnd), uintptr(wSBFlags), uintptr(bShow))
	return win.BOOL(r0)
}
