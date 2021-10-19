//go:build cocoa || (darwin && !gtk)
// +build cocoa darwin,!gtk

package dialog

import (
	"bitbucket.org/rj/goey/internal/cocoa"
)

func (m *SaveFile) show() (string, error) {
	retval := cocoa.SavePanel(m.owner.Window, m.filename)
	return retval, nil
}
