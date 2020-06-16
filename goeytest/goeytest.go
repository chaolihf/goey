package goeytest

import (
	"strings"
	"testing"
	"time"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
)

type fatalError struct{}

func (*fatalError) Error() string {
	return "sentinel for fatal error"
}

type Window interface {
	Close()
	Child() base.Element
	SetChild(base.Widget) error
}

func WithWindow(t *testing.T, init func() (Window, error)) (window Window, closer func()) {
	ready := make(chan Window, 1)
	done := make(chan struct{})
	quickCheck := strings.HasSuffix(t.Name(), "QuickCheck")

	go func() {
		winit := func() error {
			window, err := init()
			if err != nil {
				t.Errorf("failed to create window: %s", err)
				return (*fatalError)(nil)
			}
			if window == nil {
				t.Errorf("unexpected nil for window")
				return (*fatalError)(nil)
			}

			ready <- window
			return nil
		}

		err := loop.Run(winit)
		if err != nil {
			if _, ok := err.(*fatalError); ok {
				ready <- nil
			} else {
				t.Errorf("failed to run GUI loop: %s", err)
			}
		}
		close(done)
	}()

	window = <-ready
	if window == nil {
		t.SkipNow()
	}

	closer = func() {
		if testing.Verbose() && !testing.Short() && !quickCheck {
			time.Sleep(250 * time.Millisecond)
		}

		// Close the window
		err := loop.Do(func() error {
			window.Close()
			return nil
		})
		if err != nil {
			t.Errorf("failed to run loop.Do: %s", err)
		}

		// Wait for the GUI loop to terminate
		<-done
	}
	return window, closer
}
