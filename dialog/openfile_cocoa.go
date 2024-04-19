//go:build cocoa || (darwin && !gtk)
// +build cocoa darwin,!gtk

package dialog

import (
	"github.com/chaolihf/goey/internal/cocoa"
)

func (m *OpenFile) show() (string, error) {
	retval := cocoa.OpenPanel(m.owner.Window, m.filename)
	return retval, nil
}
