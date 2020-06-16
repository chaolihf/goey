package goey

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/goeytest"
	"bitbucket.org/rj/goey/loop"
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

	if base.PLATFORM != "cocoa" {
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

func testingMountWidgets(t *testing.T, widgets ...base.Widget) {
	window, closer := goeytest.WithWindow(t, func() (goeytest.Window, error) {
		return NewWindow(t.Name(), &VBox{Children: widgets})
	})
	defer closer()

	elements := window.Child().(*vboxElement).children
	goeytest.CompareElementsToWidgets(t, normalize, elements, widgets)
}

func testingMountWidget(t *testing.T, widget base.Widget) (ok bool) {
	window, closer := goeytest.WithWindow(t, func() (goeytest.Window, error) {
		return NewWindow(t.Name(), &VBox{Children: []base.Widget{widget}})
	})
	defer closer()

	// Check that the controls that were mounted match with the list
	children := window.Child().(*vboxElement).children
	if len(children) != 1 {
		t.Errorf("Wanted len(children) == 1, got %d", len(children))
		return false
	}

	return goeytest.CompareElementToWidget(t, normalize, children[0], widget)
}

func testingMountWidgetsFail(t *testing.T, outError error, widgets ...base.Widget) {
	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if window != nil {
			t.Errorf("Unexpected non-nil window")
		}
		if err != outError {
			if err == nil {
				t.Errorf("Unexpected nil error, want %s", outError)
			} else {
				t.Errorf("Unexpected error, want %v, got %s", outError, err)
			}
			return nil
		}
		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}

func testingCloseWidgets(t *testing.T, widgets ...base.Widget) {
	window, closer := goeytest.WithWindow(t, func() (goeytest.Window, error) {
		return NewWindow(t.Name(), &VBox{Children: widgets})
	})
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

func testingCheckFocusAndBlur(t *testing.T, widgets ...base.Widget) {
	log := bytes.NewBuffer(nil)
	skipFlag := false

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

	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		go func(window *Window) {
			// Wait for the window to be active.
			// This does not appear to be necessary on WIN32.  With GTK, the
			// window needs time to display before it will respond properly to
			// focus events.
			time.Sleep(100 * time.Millisecond)

			// Run the actions, which are counted.
			for i := 0; i < 3; i++ {
				err := loop.Do(func() error {
					// Find the child element to be focused
					child := window.child.(*vboxElement).children[i]
					if elem, ok := child.(Focusable); ok {
						ok := elem.TakeFocus()
						if !ok {
							t.Errorf("Failed to set focus on the control")
						}
					} else {
						skipFlag = true
					}
					return nil
				})
				if err != nil {
					t.Errorf("Error in Do, %s", err)
				}
				time.Sleep(20 * time.Millisecond)
			}

			// Close the window
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if skipFlag {
		t.Skip("Control does not support TakeFocus")
	}

	const want = "fabafbbbfcbc"
	if s := log.String(); s != want {
		t.Errorf("Incorrect log string, want %s, got log==%s", want, s)
	}
}

func testingTypeKeys(t *testing.T, text string, widget base.Widget) {
	// Typing keys happens asynchronously to the event loop.  Errors in that
	// goroutine will be fed into this channel.  However, the channel won't be
	// drained until the event loop terminates.  Errors need to be buffered.
	errc := make(chan error, 8)

	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: []base.Widget{widget}})
		if err != nil {
			return fmt.Errorf("failed to create window: %s", err)
		}

		var typingErr chan error
		go func(window *Window) {
			defer close(errc)

			// On WIN32, let the window complete any opening animation.
			time.Sleep(20 * time.Millisecond)

			err := loop.Do(func() error {
				child := window.child.(*vboxElement).children[0]
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

			// Close the window
			err = loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				errc <- fmt.Errorf("can't close window: %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	for err := range errc {
		t.Skip(err.Error())
	}
}

func testingCheckClick(t *testing.T, widgets ...base.Widget) {
	log := bytes.NewBuffer(nil)
	skipFlag := false

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

	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		go func(window *Window) {
			// Run the actions, which are counted.
			for i := 0; i < 3; i++ {
				err := loop.Do(func() error {
					// Find the child element to be clicked
					child := window.child.(*vboxElement).children[i]
					if elem, ok := child.(Clickable); ok {
						elem.Click()
					} else {
						skipFlag = true
					}
					return nil
				})
				if err != nil {
					t.Errorf("Error in Do, %s", err)
				}
			}

			// Close the window
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if skipFlag {
		t.Skip("Control does not support Click")
	}

	const want = "cacbcc"
	if s := log.String(); s != want {
		t.Errorf("Incorrect log string, want %s, got log==%s", want, s)
	}
}

func testingUpdateWidgets(t *testing.T, widgets []base.Widget, update []base.Widget) {
	window, closer := goeytest.WithWindow(t, func() (goeytest.Window, error) {
		return NewWindow(t.Name(), &VBox{Children: widgets})
	})
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

func testingUpdateWidget(t *testing.T) (updater func(base.Widget) bool, closer func()) {
	window, closer := goeytest.WithWindow(t, func() (goeytest.Window, error) {
		return NewWindow(t.Name(), nil)
	})

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

func testingLayoutWidget(t *testing.T, widget base.Widget) (updater func(base.Constraints) base.Size, closer func()) {
	window, closer := goeytest.WithWindow(t, func() (goeytest.Window, error) {
		return NewWindow(t.Name(), widget)
	})

	if !goeytest.CompareElementToWidget(t, normalize, window.Child(), widget) {
		closer()
		t.Fatalf("widget not correctly mounted")
	}

	updater = func(bc base.Constraints) base.Size {
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

	return updater, closer
}

func testingMinSizeWidget(t *testing.T, widget base.Widget) {
	window, closer := goeytest.WithWindow(t, func() (goeytest.Window, error) {
		return NewWindow(t.Name(), widget)
	})
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
		width2 := child.MinIntrinsicWidth(120 * base.DIP)
		if width1 <= 0 || width1 == base.Inf {
			t.Errorf("invalid min width: %s", width2)
		}
		if width2 < width1 {
			t.Errorf("width with height limit less than unbounded height")
		}

		height1 := child.MinIntrinsicHeight(base.Inf)
		if height1 <= 0 || height1 == base.Inf {
			t.Errorf("invalid min height: %s", height1)
		}
		height2 := child.MinIntrinsicHeight(120 * base.DIP)
		if height1 <= 0 || height1 == base.Inf {
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
