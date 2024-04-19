//go:build gtk || (linux && !cocoa) || (freebsd && !cocoa) || (openbsd && !cocoa)
// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/internal/gtk"
)

type hrElement struct {
	Control
}

func (w *HR) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountHR(parent.Handle)

	retval := &hrElement{
		Control: Control{control},
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}
