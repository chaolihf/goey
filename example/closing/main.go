// This package provides an example application built using the goey package
// that demonstrates using the OnClosing callback for windows.  Trying to close
// the window using the normal method will fail, but the button within the
// window can be used.
package main

import (
	"fmt"

	"github.com/chaolihf/goey"
	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/loop"
	"github.com/chaolihf/goey/windows"
)

var (
	mainWindow *windows.Window
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := windows.NewWindow("Closing", render())
	if err != nil {
		return err
	}
	mw.SetOnClosing(func() bool {
		// Block closing of the window
		return true
	})
	mainWindow = mw

	return nil
}

func render() base.Widget {
	return &goey.Padding{
		Insets: goey.UniformInsets(36 * goey.DIP),
		Child: &goey.Align{
			Child: &goey.Button{Text: "Close app", OnClick: func() {
				mainWindow.Close()
			}},
		},
	}
}
