package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.MustLoadDLL("kernel32.dll")
	moduser32   = syscall.MustLoadDLL("user32.dll")
	modcomctl32 = syscall.MustLoadDLL("comctl32")

	procGetDesktopWindow    = moduser32.MustFindProc("GetDesktopWindow")
	procGetWindowText       = moduser32.MustFindProc("GetWindowTextW")
	procGetWindowTextLength = moduser32.MustFindProc("GetWindowTextLengthW")
	procSetWindowText       = moduser32.MustFindProc("SetWindowTextW")
	procShowScrollBar       = moduser32.MustFindProc("ShowScrollBar")
	procInitCommonControls  = modcomctl32.MustFindProc("InitCommonControls")
)

const (
	STM_SETIMAGE = 0x0172
)

func init() {
	InitCommonControls()
}

func GetDesktopWindow() win.HWND {
	r1, _, err := syscall.Syscall(procGetDesktopWindow.Addr(), 0, 0, 0, 0)
	if err != 0 {
		panic(err)
	}
	return win.HWND(r1)
}

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

func GetWindowTextLength(hWnd win.HWND) int32 {
	r0, _, _ := syscall.Syscall(procGetWindowTextLength.Addr(), 1, uintptr(hWnd), 0, 0)
	return int32(r0)
}

func SetWindowText(hWnd win.HWND, text *uint16) win.BOOL {
	r0, _, _ := syscall.Syscall(procSetWindowText.Addr(), 2, uintptr(hWnd), uintptr(unsafe.Pointer(text)), 0)
	return win.BOOL(r0)
}

func ShowScrollBar(hWnd win.HWND, wSBFlags uint, bShow win.BOOL) win.BOOL {
	r0, _, _ := syscall.Syscall(procShowScrollBar.Addr(), 3, uintptr(hWnd), uintptr(wSBFlags), uintptr(bShow))
	return win.BOOL(r0)
}

func InitCommonControls() {
	syscall.Syscall(procInitCommonControls.Addr(), 0, 0, 0, 0)
}
