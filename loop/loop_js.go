package loop

import (
	"testing"

	"bitbucket.org/rj/goey/internal/nopanic"
)

const (
	// Flag to control behaviour of UnlockOSThread in Run.
	unlockThreadAfterRun = true
)

var (
	actions chan func()
	quit    chan struct{}
)

func initRun() error {
	// Do nothing
	return nil
}

func terminateRun() {
	// Do nothing
}

func run() {
	actions = make(chan func())
	quit = make(chan struct{})
	defer func() {
		actions = nil
		quit = nil
	}()

	ok := true
	for ok {
		select {
		case action := <-actions:
			action()
		case _, _ = <-quit:
			ok = false
		}
	}
}

func runTesting(func() error) error {
	panic("unreachable")
}

func do(action func() error) error {
	// Make channel for the return value of the action.
	err := make(chan error, 1)

	// Make the function to execute on the GUI thread.  The action needs
	// to be wrapped to transport any panics across the channel.
	actions <- func() {
		err <- nopanic.Wrap(action)
	}

	// Block on completion of action.
	return nopanic.Unwrap(<-err)
}

func stop() {
	close(quit)
}

func testMain(m *testing.M) int {
	// We need to be locked to a thread, but not to a particular
	// thread.  No need for special coordination.
	return m.Run()
}
