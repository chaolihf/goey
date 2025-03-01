package dialog

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/chaolihf/goey/loop"
	"github.com/chaolihf/win"
)

// Owner holds a pointer to the owning window.
// This type varies between platforms.
type Owner struct {
	HWnd win.HWND
}

func asyncTypeKeys(text string, initialWait time.Duration) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		time.Sleep(initialWait)
		for _, r := range text {
			inp := [2]win.KEYBD_INPUT{
				{Type: win.INPUT_KEYBOARD, Ki: win.KEYBDINPUT{}},
				{Type: win.INPUT_KEYBOARD, Ki: win.KEYBDINPUT{}},
			}

			if r == '\n' {
				inp[0].Ki.WVk = win.VK_RETURN
				inp[1].Ki.WVk = win.VK_RETURN
				inp[1].Ki.DwFlags = win.KEYEVENTF_KEYUP
			} else if r == 0x1b {
				inp[0].Ki.WVk = win.VK_ESCAPE
				inp[1].Ki.WVk = win.VK_ESCAPE
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
					return fmt.Errorf("windows error, %x", win.GetLastError())
				}
				return nil
			})
			if err != nil {
				errs <- err
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	return errs
}
