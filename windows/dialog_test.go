//go:build !js
// +build !js

package windows_test

import (
	"fmt"

	"bitbucket.org/rj/goey/loop"
	"bitbucket.org/rj/goey/windows"
)

func ExampleWindow_Message() {
	// All calls that modify GUI objects need to be schedule ont he GUI thread.
	// This callback will be used to create the top-level window.
	createWindow := func() error {
		// Create a top-level window.
		mw, err := windows.NewWindow("Test", nil /*empty window*/)
		if err != nil {
			// This error will be reported back up through the call to
			// Run below.  No need to print or log it here.
			return err
		}

		// We can start a goroutine, but note that we can't modify GUI objects
		// directly.
		go func() {
			// Show the error message.
			_ = loop.Do(func() error {
				return mw.Message("This is an example message.").WithInfo().Show()
			})

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			_ = loop.Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	// Start the GUI thread.
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
