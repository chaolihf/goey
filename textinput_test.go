package goey

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/windows"
)

func ExampleTextInput() {
	// In a full application, this variable would be updated to point to
	// the main window for the application.
	var mainWindow *windows.Window

	// These functions are used to update the GUI.  See below.
	var update func()
	var render func() base.Widget

	// Update function
	update = func() {
		err := mainWindow.SetChild(render())
		if err != nil {
			panic(err)
		}
	}

	// Render function generates a tree of Widgets to describe the desired
	// state of the GUI.
	render = func() base.Widget {
		// Prep - text for the button
		// The GUI contains a single widget, this button.
		return &VBox{Children: []base.Widget{
			&Label{Text: "Enter you text below:"},
			&TextInput{
				Value:       "",
				Placeholder: "Enter your data here",
				OnChange: func(value string) {
					fmt.Println("Change: ", value)
					// In a real example, you would update your data, and then
					// need to render the window again.
					update()
				},
			},
		}}
	}
}

func textinputValues(values []reflect.Value, rand *rand.Rand) {
	// Get a string
	labelValues(values, rand)

	// Create a choices for disabled and default
	values[1] = reflect.ValueOf(rand.Uint64()%2 == 0)
	values[2] = reflect.ValueOf(rand.Uint64()%2 == 0)
	values[3] = reflect.ValueOf(rand.Uint64()%2 == 0)
}

func TestTextInputMount(t *testing.T) {
	testMountWidgets(t,
		&TextInput{Value: "A"},
		&TextInput{Value: "B", Placeholder: "..."},
		&TextInput{Value: "C", Disabled: true},
		&TextInput{Value: "D", ReadOnly: true},
		&TextInput{Value: "E", Password: true},
	)

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		f := func(value string, disabled, password, readonly bool) bool {
			return checkMountWidget(t, &TextInput{Value: value, Disabled: disabled, Password: password, ReadOnly: readonly})
		}
		if err := quick.Check(f, &quick.Config{Values: textinputValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}

func TestTextInputClose(t *testing.T) {
	testCloseWidgets(t,
		&TextInput{Value: "A"},
		&TextInput{Value: "B", Placeholder: "..."},
		&TextInput{Value: "C", Disabled: true},
		&TextInput{Value: "D", ReadOnly: true},
		&TextInput{Value: "E", Password: true},
	)
}

func TestTextInputOnFocus(t *testing.T) {
	testCheckFocusAndBlur(t,
		&TextInput{},
		&TextInput{},
		&TextInput{},
	)

	// On some platforms, the password control is a separate type, and so may
	// may have a parallel implementation.
	testCheckFocusAndBlur(t,
		&TextInput{Password: true},
		&TextInput{Password: true},
		&TextInput{Password: true},
	)
}

func TestTextInputOnChange(t *testing.T) {
	log := bytes.NewBuffer(nil)

	testTypeKeys(t, "Hello",
		&TextInput{OnChange: func(v string) {
			log.WriteString(v)
			log.WriteString("\x1E")
		}})

	const want = "H\x1EHe\x1EHel\x1EHell\x1EHello\x1E"
	if got := log.String(); got != want {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func TestTextInputOnEnterKey(t *testing.T) {
	log := bytes.NewBuffer(nil)

	testTypeKeys(t, "Hello\n",
		&TextInput{OnEnterKey: func(v string) {
			log.WriteString(v)
		}})

	const want = "Hello"
	if got := log.String(); got != want {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func TestTextInputUpdateProps(t *testing.T) {
	testUpdateWidgets(t, []base.Widget{
		&TextInput{Value: "A"},
		&TextInput{Value: "B", Placeholder: "..."},
		&TextInput{Value: "C", Disabled: true},
		&TextInput{Value: "D", ReadOnly: true},
	}, []base.Widget{
		&TextInput{Value: "AA", ReadOnly: true},
		&TextInput{Value: "BA", Disabled: true},
		&TextInput{Value: "CA", Placeholder: "***", Disabled: false},
		&TextInput{Value: "DA"},
	})

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		updater, closer := checkUpdateWidget(t)
		defer closer()

		f := func(value string, disabled, password, readonly bool) bool {
			// Deliberately ignoring flag for password.
			// Changing text input into a password field is not supported on
			// all platforms.
			_ = password
			return updater(&TextInput{Value: value, Disabled: disabled, ReadOnly: readonly})
		}
		if err := quick.Check(f, &quick.Config{Values: textinputValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}

func TestTextInputLayout(t *testing.T) {
	testLayoutWidget(t, &TextInput{Value: "AB"})
}

func TestTextInputMinSize(t *testing.T) {
	testMinSizeWidget(t, &TextInput{Value: "AB"})
}
