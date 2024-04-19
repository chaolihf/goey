//go:build cocoa || (darwin && !gtk)
// +build cocoa darwin,!gtk

package dialog

import (
	"github.com/chaolihf/goey/internal/cocoa"
)

func (m *Message) show() error {
	cocoa.MessageDialog(m.owner.Window, m.text, m.title, byte(m.icon))
	return nil
}

func (m *Message) withError() {
	m.icon = 'e'
}

func (m *Message) withWarn() {
	m.icon = 'w'
}

func (m *Message) withInfo() {
	m.icon = 'i'
}
