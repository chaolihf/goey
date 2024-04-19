package goeytest

import (
	"reflect"
	"testing"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/loop"
)

// A Proper is an element that can extract it's properties from the underlying
// GUI control, and recreate the matching base.Widget.
type Proper interface {
	Props() base.Widget
}

func equal(t *testing.T, normalize func(*testing.T, base.Widget), lhs, rhs base.Widget) bool {
	// Normalize (or canonicalize) the props used to construct the element.
	normalize(t, rhs)
	// Compare the widgets' properties.
	return reflect.DeepEqual(lhs, rhs)
}

// CompareElementToWidget returns true if the element and widget are equal.
// The element must satisfy the Proper interface.
func CompareElementToWidget(t *testing.T, normalize func(*testing.T, base.Widget), element base.Element, widget base.Widget) bool {
	return element.Kind() == widget.Kind() &&
		equal(t, normalize, element.(Proper).Props(), widget)
}

// CompareElementsToWidgets returns true if all of the the elements in a slice are equal to the correspondings widget.
// All of the elements must satisfy the Proper interface.
func CompareElementsToWidgets(t *testing.T, normalize func(*testing.T, base.Widget), elements []base.Element, widgets []base.Widget) {
	if len(elements) != len(widgets) {
		t.Errorf("wanted len(elements) == len(widgets), got %d and %d",
			len(elements),
			len(widgets),
		)
		return
	}

	for i := range elements {
		if n1, n2 := elements[i].Kind(), widgets[i].Kind(); n1 != n2 {
			t.Errorf("Wanted children[%d].Kind() == widgets[%d].Kind(), got %s and %s", i, i, n1, n2)
		} else if widget, ok := elements[i].(Proper); ok {
			var data base.Widget
			err := loop.Do(func() error {
				data = widget.Props()
				return nil
			})
			if err != nil {
				t.Fatalf("error in loop.Do: %s", err)
			}

			if n1, n2 := data.Kind(), widgets[i].Kind(); n1 != n2 {
				t.Errorf("Wanted data.Kind() == widgets[%d].Kind(), got %s, want %s", i, n1, n2)
			}
			if !equal(t, normalize, data, widgets[i]) {
				t.Errorf("Wanted data == widgets[%d], got %v, want %v", i, data, widgets[i])
			}
		} else {
			t.Skipf("Cannot verify props of child %d", i)
		}
	}
}
