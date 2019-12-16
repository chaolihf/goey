// +build cocoa darwin,!gtk

package dialog

import (
	"bitbucket.org/rj/goey/internal/cocoa"
	"time"
)

type dialogImpl struct {
	parent *cocoa.Window
}

func asyncTypeKeys(text string, initialWait time.Duration) chan error {
	err := make(chan error, 1)

	go func() {
		defer close(err)

		// TODO
	}()

	return err
}
