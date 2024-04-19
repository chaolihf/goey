package goeytest

import (
	"strings"
	"testing"
	"time"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/loop"
	"github.com/chaolihf/goey/windows"
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

// WithWindow initializes a window and a GUI event loop that can be used to test
// widgets.  When testing is complete, callers should use the return callback to
// close the window and terminate the event loop.
//
// Depending on platform, there may be a period between when the window is
// created and when it is mapped (GTK), or when animations have completed
// (windows).  During this time, injecting events such as button presses and key
// presses may fail.
func WithWindow(t *testing.T, widget base.Widget) (window *windows.Window, closer func()) {
	ready := make(chan *windows.Window, 1)
	done := make(chan struct{})
	quickCheck := strings.HasSuffix(t.Name(), "QuickCheck")

	go func() {
		winit := func() error {
			window, err := windows.NewWindow(t.Name(), widget)
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
