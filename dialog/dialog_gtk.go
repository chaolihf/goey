//go:build gtk || (linux && !cocoa) || (freebsd && !cocoa) || (openbsd && !cocoa)
// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package dialog

import (
	"time"

	"bitbucket.org/rj/goey/internal/gtk"
	"bitbucket.org/rj/goey/loop"
)

// Owner holds a pointer to the owning window.
// This type varies between platforms.
type Owner struct {
	Handle uintptr
}

var (
	activeDialogForTesting uintptr
)

func asyncTypeKeys(text string, initialWait time.Duration) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		time.Sleep(initialWait)
		for _, r := range text {
			err := loop.Do(func() error {
				if activeDialogForTesting == 0 {
					panic("dialog is closed")
				}
				gtk.WidgetSendKey(activeDialogForTesting, uint(r), false)
				return nil
			})
			if err != nil {
				errs <- err
				return
			}
			time.Sleep(50 * time.Millisecond)

			err = loop.Do(func() error {
				if activeDialogForTesting != 0 {
					gtk.WidgetSendKey(activeDialogForTesting, uint(r), true)
				}
				return nil
			})
			if err != nil {
				errs <- err
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	return errs
}
