package goey

import (
	"time"
)

var (
	dateInputKind = Kind{"dateinput"}
)

// DateInput describes a widget that users input or update a single date.
// The model for the value is a string value.
type DateInput struct {
	Value    time.Time             // Values is the current string for the field
	Disabled bool                  // Disabled is a flag indicating that the user cannot interact with this field
	OnChange func(value time.Time) // OnChange will be called whenever the user changes the value for this field
	OnFocus  func()                // OnFocus will be called whenever the field receives the keyboard focus
	OnBlur   func()                // OnBlur will be called whenever the field loses the keyboard focus
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*DateInput) Kind() *Kind {
	return &dateInputKind
}

// Mount creates a text field in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *DateInput) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedDateInput) Kind() *Kind {
	return &dateInputKind
}

func (w *mountedDateInput) UpdateProps(data Widget) error {
	return w.updateProps(data.(*DateInput))
}
