//go:build !js
// +build !js

package windows

// Dialogs are currently not supported on the JS platform.

import (
	"github.com/chaolihf/goey/dialog"
)

// Message returns a builder that can be used to construct a message
// dialog, and then show that dialog.
func (w *Window) Message(text string) *dialog.Message {
	ret := dialog.NewMessage(text)
	w.message(ret)
	return ret
}

// OpenFileDialog returns a builder that can be used to construct an open file
// dialog, and then show that dialog.
func (w *Window) OpenFileDialog() *dialog.OpenFile {
	ret := dialog.NewOpenFile()
	w.openfiledialog(ret)
	return ret
}

// SaveFileDialog returns a builder that can be used to construct a save file
// dialog, and then show that dialog.
func (w *Window) SaveFileDialog() *dialog.SaveFile {
	ret := dialog.NewSaveFile()
	w.savefiledialog(ret)
	return ret
}
