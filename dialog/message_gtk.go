//go:build gtk || (linux && !cocoa) || (freebsd && !cocoa) || (openbsd && !cocoa)
// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package dialog

import (
	"github.com/chaolihf/goey/internal/gtk"
)

func (m *Message) show() error {
	dlg := gtk.MountMessageDialog(m.owner.Handle, m.title, m.icon, m.text)
	activeDialogForTesting = dlg
	defer func() {
		activeDialogForTesting = 0
		gtk.WidgetClose(dlg)
	}()

	gtk.DialogRun(dlg)
	return nil
}

func (m *Message) withError() {
	m.icon = gtk.MessageDialogWithError()
}

func (m *Message) withWarn() {
	m.icon = gtk.MessageDialogWithWarn()
}

func (m *Message) withInfo() {
	m.icon = gtk.MessageDialogWithInfo()
}
