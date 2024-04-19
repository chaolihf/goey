package goey

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/goeytest"
	"github.com/chaolihf/goey/loop"
	"github.com/chaolihf/goey/windows"
)

type Proper interface {
	Props() base.Widget
}

type Clickable interface {
	Click()
}

type Focusable interface {
	TakeFocus() bool
}

type Typeable interface {
	TakeFocus() bool
	TypeKeys(text string) chan error
}

func TestMain(m *testing.M) {
	// On Cocoa, the GUI even
	loop.TestMain(m)
}

func normalize(t *testing.T, rhs base.Widget) {
	if base.PLATFORM == "windows" {
		// On windows, the message EM_GETCUEBANNER does not work unless the manifest
		// is set correctly.  This cannot be done for the package, since that
		// manifest will conflict with the manifest of any app.
		if value := reflect.ValueOf(rhs).Elem().FieldByName("Placeholder"); value.IsValid() {
			placeholder := value.String()
			if placeholder != "" {
				t.Logf("Zeroing 'Placeholder' field during test")
			}
			value.SetString("")
		}
	} else if base.PLATFORM == "gtk" {
		// With GTK, this package is using a GtkTextView to create
		// the multi-line text editor, and that widget does not support
		// placeholders.
		if elem, ok := rhs.(*TextArea); ok {
			if elem.Placeholder != "" {
				t.Logf("Zeroing 'Placeholder' field during test")
			}
			elem.Placeholder = ""
		}
	}

	if base.PLATFORM == "windows" || base.PLATFORM == "gtk" {
		// On both windows and GTK, the props method only return RGBA images.
		if value := reflect.ValueOf(rhs).Elem().FieldByName("Image"); value.IsValid() {
			if prop, ok := value.Interface().(*image.Gray); ok {
				t.Logf("Converting 'Image' field from *image.Gray to *image.RGBA")
				bounds := prop.Bounds()
				img := image.NewRGBA(bounds)
				draw.Draw(img, bounds, prop, bounds.Min, draw.Src)
				value.Set(reflect.ValueOf(img))
			}
		}
	}

	if value := reflect.ValueOf(rhs).Elem().FieldByName("Child"); value.IsValid() {
		if child := value.Interface(); child != nil {
			normalize(t, child.(base.Widget))
		}
	}
}

func testMountWidgets(t *testing.T, widgets ...base.Widget) {
	window, closer := goeytest.WithWindow(t, &VBox{Children: widgets})
	defer closer()

	elements := window.Child().(*vboxElement).children
	goeytest.CompareElementsToWidgets(t, normalize, elements, widgets)
}

func checkMountWidget(t *testing.T, widget base.Widget) (ok bool) {
	window, closer := goeytest.WithWindow(t, &VBox{Children: []base.Widget{widget}})
	defer closer()

	// Check that the controls that were mounted match with the list
	children := window.Child().(*vboxElement).children
	if len(children) != 1 {
		t.Errorf("Wanted len(children) == 1, got %d", len(children))
		return false
	}

	return goeytest.CompareElementToWidget(t, normalize, children[0], widget)
}

func testMountWidgetsFail(t *testing.T, outError error, widgets ...base.Widget) {
	init := func() error {
		window, err := windows.NewWindow(t.Name(), &VBox{Children: widgets})
		if window != nil {
			t.Errorf("unexpected non-nil window")
		}
		if err != outError {
			if err == nil {
				t.Errorf("unexpected nil error: want %s", outError)
			} else {
				t.Errorf("unexpected error: want %v, got %s", outError, err)
			}
			return nil
		}
		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("failed to run GUI loop: %s", err)
	}
}

func testCloseWidgets(t *testing.T, widgets ...base.Widget) {
	window, closer := goeytest.WithWindow(t, &VBox{Children: widgets})
	defer closer()

	elements := window.Child().(*vboxElement).children
	goeytest.CompareElementsToWidgets(t, normalize, elements, widgets)

	err := loop.Do(func() error {
		return window.SetChild(nil)
	})
	if err != nil {
		t.Fatalf("error in loop.Do: %s", err)
	}
}

func testCheckFocusAndBlur(t *testing.T, widgets ...base.Widget) {
	// Rewrite the widgets to add the OnFocus and OnBlur event handlers.
	log := bytes.NewBuffer(nil)
	for i := byte(0); i < 3; i++ {
		s := reflect.ValueOf(widgets[i])
		letter := 'a' + i
		s.Elem().FieldByName("OnFocus").Set(reflect.ValueOf(func() {
			log.Write([]byte{'f', letter})
		}))
		s.Elem().FieldByName("OnBlur").Set(reflect.ValueOf(func() {
			log.Write([]byte{'b', letter})
		}))
	}

	func() {
		window, closer := goeytest.WithWindow(t, &VBox{Children: widgets})
		defer closer()

		// Wait for the window to be active.
		// This does not appear to be necessary on WIN32.  With GTK, the
		// window needs time to display before it will respond properly to
		// focus events.
		time.Sleep(100 * time.Millisecond)

		// Run the actions, which are counted.
		for i := 0; i < 3; i++ {
			child := window.Child().(*vboxElement).children[i]
			// Find the child element to be focused
			if elem, ok := child.(Focusable); ok {
				err := loop.Do(func() error {
					ok := elem.TakeFocus()
					if !ok {
						t.Errorf("Failed to set focus on the control")
					}
					return nil
				})
				if err != nil {
					t.Errorf("error in loop.Do: %s", err)
				}
			} else {
				t.Skip("control does not support method TakeFocus")
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()

	const want = "fabafbbbfcbc"
	if s := log.String(); s != want {
		t.Errorf("Incorrect log string, want %s, got log==%s", want, s)
	}
}

func testTypeKeys(t *testing.T, text string, widget base.Widget) {
	// Typing keys happens asynchronously to the event loop.  Errors in that
	// goroutine will be fed into this channel.  However, the channel won't be
	// drained until the event loop terminates.  Errors need to be buffered.
	errc := make(chan error, 8)

	window, closer := goeytest.WithWindow(t, &VBox{Children: []base.Widget{widget}})
	defer closer()

	var typingErr chan error
	go func(window *windows.Window) {
		defer close(errc)

		// On WIN32, let the window complete any opening animation.
		time.Sleep(20 * time.Millisecond)

		err := loop.Do(func() error {
			child := window.Child().(*vboxElement).children[0]
			if elem, ok := child.(Typeable); ok {
				if elem.TakeFocus() {
					typingErr = elem.TypeKeys(text)
				} else {
					return fmt.Errorf("control failed to take focus")
				}
			} else {
				return fmt.Errorf("control does not support method TypeKeys")
			}
			return nil
		})
		if err != nil {
			errc <- err
		}

		// Wait for typing to complete, and check for errors
		if typingErr != nil {
			for v := range typingErr {
				errc <- fmt.Errorf("failed to type keys: %s", v)
			}
		}
	}(window)

	for err := range errc {
		t.Skip(err.Error())
	}
}

func testCheckClick(t *testing.T, widgets ...base.Widget) {
	// Rewrite the widgets to add the onclick events
	log := bytes.NewBuffer(nil)
	for i := byte(0); i < 3; i++ {
		letter := 'a' + i
		if elem, ok := widgets[i].(*Checkbox); ok {
			// Chain the onclick callback for the element.
			chainCallback := elem.OnChange
			// Add wrapper to write to the test log.
			elem.OnChange = func(value bool) {
				log.Write([]byte{'c', letter})
				chainCallback(value)
			}
		} else {
			s := reflect.ValueOf(widgets[i])
			s.Elem().FieldByName("OnClick").Set(reflect.ValueOf(func() {
				log.Write([]byte{'c', letter})
			}))
		}
	}

	window, closer := goeytest.WithWindow(t, &VBox{Children: widgets})
	defer closer()

	// Run the actions, which are counted.
	for i := 0; i < 3; i++ {
		// Find the child element to be clicked
		child := window.Child().(*vboxElement).children[i]

		if elem, ok := child.(Clickable); ok {
			err := loop.Do(func() error {
				elem.Click()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		} else {
			t.Skip("control does not support method Click")
		}
	}

	const want = "cacbcc"
	if s := log.String(); s != want {
		t.Errorf("Incorrect log string, want %s, got log==%s", want, s)
	}
}

func testUpdateWidgets(t *testing.T, widgets []base.Widget, update []base.Widget) {
	window, closer := goeytest.WithWindow(t, &VBox{Children: widgets})
	defer closer()

	elements := window.Child().(*vboxElement).children
	goeytest.CompareElementsToWidgets(t, normalize, elements, widgets)

	err := loop.Do(func() error {
		return window.SetChild(&VBox{Children: update})
	})
	if err != nil {
		t.Fatalf("error in loop.Do: %s", err)
	}

	elements = window.Child().(*vboxElement).children
	goeytest.CompareElementsToWidgets(t, normalize, elements, update)
}

func checkUpdateWidget(t *testing.T) (updater func(base.Widget) bool, closer func()) {
	window, closer := goeytest.WithWindow(t, nil)

	updater = func(w base.Widget) bool {
		err := loop.Do(func() error {
			return window.SetChild(w)
		})
		if err != nil {
			t.Fatalf("error in loop.Do: %s", err)
			return false
		}

		return goeytest.CompareElementToWidget(t, normalize, window.Child(), w)
	}

	return updater, closer
}

func testLayoutWidget(t *testing.T, widget base.Widget) {
	window, closer := goeytest.WithWindow(t, widget)
	defer closer()

	if !goeytest.CompareElementToWidget(t, normalize, window.Child(), widget) {
		closer()
		t.Fatalf("widget not correctly mounted")
	}

	updater := func(bc base.Constraints) base.Size {
		size := base.Size{}
		err := loop.Do(func() error {
			size = window.Child().Layout(bc)
			return nil
		})
		if err != nil {
			t.Fatalf("error in loop.Do: %s", err)
		}
		return size
	}

	cases := []struct {
		name string
		bc   base.Constraints
	}{
		{"unconstrainied", base.Unconstrained()},
		{"expand-with-min", base.Constraints{base.Size{base.Inf / 2, base.Inf / 2}, base.Size{base.Inf, base.Inf}}},
		{"loose", base.Loose(base.Size{96 * DIP, 24 * DIP})},
		{"loose-height", base.LooseHeight(24 * DIP)},
		{"loose-width", base.LooseWidth(96 * DIP)},
		{"tight", base.Tight(base.Size{96 * DIP, 24 * DIP})},
		{"tight-height", base.TightHeight(24 * DIP)},
		{"tight-width", base.TightWidth(96 * DIP)},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			size := updater(v.bc)
			if !v.bc.IsSatisfiedBy(size) {
				t.Errorf("layout does not respect constraints")
			}
		})
	}
}

func testMinSizeWidget(t *testing.T, widget base.Widget) {
	window, closer := goeytest.WithWindow(t, widget)
	defer closer()

	child := window.Child()
	if !goeytest.CompareElementToWidget(t, normalize, child, widget) {
		t.Errorf("widget not correctly mounted")
		return
	}

	err := loop.Do(func() error {
		width1 := child.MinIntrinsicWidth(base.Inf)
		if width1 <= 0 || width1 == base.Inf {
			t.Errorf("invalid min width: %s", width1)
		}
		t.Logf("MinIntrinsicWidth(Inf): %s", width1)

		width2 := child.MinIntrinsicWidth(120 * base.DIP)
		if width2 <= 0 || width2 == base.Inf {
			t.Errorf("invalid min width: %s", width2)
		}
		if width2 < width1 {
			t.Errorf("width with height limit less than unbounded height")
		}

		height1 := child.MinIntrinsicHeight(base.Inf)
		if height1 <= 0 || height1 == base.Inf {
			t.Errorf("invalid min height: %s", height1)
		}
		t.Logf("MinIntrinsicHeight(Inf): %s", height1)

		height2 := child.MinIntrinsicHeight(120 * base.DIP)
		if height2 <= 0 || height2 == base.Inf {
			t.Errorf("invalid min height: %s", height2)
		}
		if height2 < height1 {
			t.Errorf("height with width limit less than unbounded width")
		}

		return nil
	})
	if err != nil {
		t.Fatalf("error in loop.Do: %s", err)
	}
}
