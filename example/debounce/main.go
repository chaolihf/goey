// This package provides an example application built using the goey package
// that shows how an event be debounced.  Events from the user typing into a
// field are held until the user stops typing.
package main

import (
	"fmt"
	"os"

	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
	"bitbucket.org/rj/goey/windows"
)

var (
	mainWindow *windows.Window
	username   string
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
	mw, err := windows.NewWindow("Debounce", render())
	if err != nil {
		return err
	}

	// We store a copy of the pointer to the window so that we can update the
	// GUI at a later time.
	mainWindow = mw

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
	// The greating will change depending on what the user has entered as
	// their name.
	text := "To whom am I speaking?"
	if username != "" {
		text = "Hello, " + username + "!"
	}

	// We return a widget describing the desired state of the GUI.  Note that
	// this is data only, and no changes have been effected yet.
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			Children: []base.Widget{
				&goey.Label{Text: "Your name, please?"},
				&goey.TextInput{
					Value:       username,
					Placeholder: "Your name",
					OnChange: Debounce(func(s string) {
						username = s
						updateWindow()
					}),
				},
				&goey.HR{},
				&goey.P{
					Text: text,
				},
			},
		},
	}
}
