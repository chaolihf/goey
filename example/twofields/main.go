// This package provides an example application built using the goey package
// that demonstrates two multiline text fields.  A status bar shows the combined
// count of characters in both fields, showing how a dynamic GUI can be easily
// kept in sync with changes to the application's data.
//
// This example also shows the use of the Expand widget to have some children
// of the VBox expand and consume any available vertical space.
package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/chaolihf/goey"
	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/loop"
	"github.com/chaolihf/goey/windows"
)

var (
	mainWindow *windows.Window
	text       [2]string
)

func main() {
	flag.StringVar(&text[0], "text0", "", "Initial text for the first field")
	flag.StringVar(&text[1], "text1", "", "Initial text for the second field")
	flag.Parse()

	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := windows.NewWindow("Two Fields", render())
	if err != nil {
		return err
	}
	mw.SetScroll(false, true)
	mainWindow = mw
	return nil
}

func updateWindow() {
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func render() base.Widget {
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			Children: []base.Widget{
				&goey.Label{Text: "This is the most important field:"},
				&goey.Expand{Child: &goey.TextArea{
					Value:       text[0],
					Placeholder: "You should type something here.",
					OnChange: func(value string) {
						text[0] = value
						updateWindow()
					},
					OnFocus: onfocus(1),
					OnBlur:  onblur(1),
				}},
				&goey.Label{Text: "This is a secondary field:"},
				&goey.Expand{Child: &goey.TextArea{
					Value:       text[1],
					Placeholder: "...and here.",
					OnChange: func(value string) {
						text[1] = value
						updateWindow()
					},
					OnFocus: onfocus(2),
					OnBlur:  onblur(2),
				}},
				&goey.HR{},
				&goey.Label{Text: "The total character count is:  " + strconv.Itoa(len(text[0])+len(text[1]))},
			},
		},
	}
}

func onfocus(ndx int) func() {
	return func() {
		fmt.Println("focus", ndx)
	}
}

func onblur(ndx int) func() {
	return func() {
		fmt.Println("blur", ndx)
	}
}
