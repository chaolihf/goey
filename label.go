package goey

import (
	"github.com/chaolihf/goey/base"
)

var (
	labelKind = base.NewKind("github.com/chaolihf/goey.Label")
)

// Label describes a widget that provides a descriptive label for other fields.
//
// Labels should not be empty, and should not contain leading or trailing
// spaces.  If violated, the behaviour of the Label will depend on the GUI
// platform targeted.
type Label struct {
	Text string // Text is the contents of the label
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Label) Kind() *base.Kind {
	return &labelKind
}

// Mount creates a label control in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Label) Mount(parent base.Control) (base.Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*labelElement) Kind() *base.Kind {
	return &labelKind
}

func (w *labelElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*Label))
}
