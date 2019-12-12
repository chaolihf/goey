package dialog

import (
	"time"

	"bitbucket.org/rj/goey/internal/gtk"
	"bitbucket.org/rj/goey/loop"
)

type dialogImpl struct {
	parent uintptr
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
				gtk.WidgetSendKey(activeDialogForTesting, uint(r), false)
				return nil
			})
			if err != nil {
				errs <- err
				return
			}
			time.Sleep(50 * time.Millisecond)

			err = loop.Do(func() error {
				gtk.WidgetSendKey(activeDialogForTesting, uint(r), true)
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
