// This package provides an example application built using the goey package.
// that demontrates two multiline text fields.  A status bar shows the combined
// count of characters in both fields, showing how the GUI can be easily kept
// in sync with changes to data.
package main

import (
	"fmt"
	"strconv"

	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
)

var (
	mainWindow     *goey.Window
	text           [2]string
	characterCount [2]int
	wordCount      [2]int
)

func main() {
	err := goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := goey.NewWindow("Two Fields", render())
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
				&goey.TextArea{Value: text[0], Placeholder: "You should type something here.",
					OnChange: func(value string) {
						text[0] = value
						characterCount[0] = len(value)
						updateWindow()
					},
					OnFocus: onfocus(1),
					OnBlur:  onblur(1),
				},
				&goey.Label{Text: "This is a secondary field:"},
				&goey.TextArea{Value: text[1], Placeholder: "...and here.",
					OnChange: func(value string) {
						text[1] = value
						characterCount[1] = len(value)
						updateWindow()
					},
					OnFocus: onfocus(2),
					OnBlur:  onblur(2),
				},
				&goey.HR{},
				&goey.Label{Text: "The total character count is:  " + strconv.Itoa(characterCount[0]+characterCount[1])},
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
