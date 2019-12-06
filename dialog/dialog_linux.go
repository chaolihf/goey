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

func typeKeys(text string) chan error {
	err := make(chan error, 1)

	go func() {
		defer close(err)

		time.Sleep(500 * time.Millisecond)
		for _, r := range text {
			loop.Do(func() error {
				gtk.WidgetSendKey(activeDialogForTesting, uint(r), false)
				return nil
			})
			time.Sleep(50 * time.Millisecond)

			loop.Do(func() error {
				gtk.WidgetSendKey(activeDialogForTesting, uint(r), true)
				return nil
			})
			time.Sleep(50 * time.Millisecond)
		}
	}()

	return err
}
