//go:build gtk || (linux && !cocoa) || (freebsd && !cocoa) || (openbsd && !cocoa)
// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package dialog

import (
	"github.com/chaolihf/goey/internal/gtk"
)

func (m *SaveFile) show() (string, error) {
	dlg := gtk.MountSaveDialog(m.owner.Handle, m.title, m.filename)
	activeDialogForTesting = dlg
	defer func() {
		activeDialogForTesting = 0
		gtk.WidgetClose(dlg)
	}()

	for _, v := range m.filters {
		gtk.DialogAddFilter(dlg, v.name, v.pattern)
	}

	rc := gtk.DialogRun(dlg)
	if rc != gtk.DialogResponseAccept() {
		return "", nil
	}
	return gtk.DialogGetFilename(dlg), nil
}
