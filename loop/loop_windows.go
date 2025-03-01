package loop

import (
	"sync/atomic"
	"syscall"
	"testing"
	"unsafe"

	"github.com/chaolihf/goey/internal/nopanic"
	"github.com/chaolihf/win"
)

const (
	// Flag to control behaviour of UnlockOSThread in Run.
	isOSThreadLockedAtInit = false
)

var (
	atomPost win.ATOM
	hwndPost win.HWND
	namePost = [...]uint16{'G', 'o', 'e', 'y', 'P', 'o', 's', 't', 'W', 'i', 'n', 'd', 'o', 'w', 0}

	activeWindow uintptr

	postMessageAction = make(chan func() error, 1)
	postMessageErr    = make(chan error, 1)
)

func initRun() error {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		// Not sure that the call to GetModuleHandle can ever fail when the
		// argument is nil.  The handle for the current .exe is certainly
		// valid?
		return syscall.GetLastError()
	}

	// Make sure that we have registered a class for the hidden window.
	if atomPost == 0 {
		wc := win.WNDCLASSEX{
			CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
			HInstance:     hInstance,
			LpfnWndProc:   syscall.NewCallback(postWindowProc),
			LpszClassName: &namePost[0],
		}

		atomPost = win.RegisterClassEx(&wc)
		if atomPost == 0 {
			return syscall.GetLastError()
		}
	}

	// Create the hidden window.
	hwndPost = win.CreateWindowEx(0, &namePost[0], nil, 0,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		win.HWND_DESKTOP, 0, hInstance, nil)
	if hwndPost == 0 {
		return syscall.GetLastError()
	}
	return nil
}

func terminateRun() {
	win.DestroyWindow(hwndPost)
	hwndPost = 0
}

func run() {
	// Run the message loop
	for loop() {
	}
}

func runTesting(func() error) error {
	panic("unreachable")
}

func do(action func() error) error {
	// Let the GUI thread know that an action is coming.
	win.PostMessage(hwndPost, win.WM_USER, 0, 0)

	// Send the action.
	postMessageAction <- action

	// Block for and return err.
	return nopanic.Unwrap(<-postMessageErr)
}

func loop() (ok bool) {
	// Obtain a copy of the next message from the queue.
	var msg win.MSG
	win.GetMessage(&msg, 0, 0, 0)

	// Processing for application wide messages are handled in this block.
	if msg.Message == win.WM_QUIT {
		return false
	}

	// Dispatch message.
	if !win.IsDialogMessage(win.HWND(activeWindow), &msg) {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
	return true
}

func stop() {
	win.PostQuitMessage(0)
}

func SetActiveWindow(hwnd win.HWND) {
	atomic.StoreUintptr(&activeWindow, uintptr(hwnd))
}

func postWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case win.WM_USER:
		postMessageErr <- nopanic.Wrap(<-postMessageAction)
		return 0
	}

	// Let the default window proc handle all other messages
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func testMain(m *testing.M) int {
	// On Windows, we need to be locked to a thread, but not to a particular
	// thread.  No need for special coordination.
	return m.Run()
}
