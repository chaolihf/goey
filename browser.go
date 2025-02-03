package goey

import (
	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/win"
)

var (
	browserKind = base.NewKind("github.com/chaolihf/goey.Browser")
)

type Browser struct {
	EmbedWindow win.HWND
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Browser) Kind() *base.Kind {
	return &browserKind
}

// Mount creates a browser control in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Browser) Mount(parent base.Control) (base.Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*browserElement) Kind() *base.Kind {
	return &browserKind
}

func (w *browserElement) UpdateProps(data base.Widget) error {
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Browser))
}
