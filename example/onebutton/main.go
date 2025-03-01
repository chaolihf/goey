// This package provides an example application built using the goey package
// that shows a single button.  The button is centered in the window, and, when
// the button is clicked, the button's caption is changed to keep a running
// total.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/chaolihf/goey"
	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/loop"
	"github.com/chaolihf/goey/windows"
)

var (
	mainWindow *windows.Window
	clickCount int
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func createWindow() error {
	// This is the callback used to initialize the GUI state.  For this simple
	// example, we need to create a new top-level window, and set a child
	// widget.
	mw, err := windows.NewWindow("One Button", render())
	if err != nil {
		return err
	}

	// We store a copy of the pointer to the window so that we can update the
	// GUI at a later time.
	mainWindow = mw
	mw.SetOnResize(func(width, height int) bool {
		fmt.Printf("resize window width %d,height %d\n", width, height)
		return true
	})
	return nil
}

func updateWindow() {
	// To update the window, we generate a new widget for the contents of the
	// top-level window.
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}

func render() base.Widget {
	// The text for the button will depend on how many times it has been
	// clicked.  Build the string for the button's caption.
	text := "Click me!"
	if clickCount > 0 {
		text = text + "  (" + strconv.Itoa(clickCount) + ")"
	}

	// We return a widget describing the desired state of the GUI.  Note that
	// this is data only, and no changes have been effected yet.
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.Align{
			Child: &goey.Button{Text: text, OnClick: func() {
				// Side-effect for clicking the button.
				clickCount++
				// Update the contents of the top-level window.
				updateWindow()
				showWindowSize()

			}},
		},
	}
}

func showWindowSize() {
	width, height := mainWindow.GetSize()
	fmt.Printf("resize window width %d,height %d\n", width, height)
}
