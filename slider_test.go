package goey

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/loop"
	"github.com/chaolihf/goey/windows"
)

func ExampleSlider() {
	value := 0.0

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
		text := "Value: " + strconv.FormatFloat(value, 'f', 1, 64)
		// The GUI contains a single widget, this button.
		return &VBox{
			AlignMain:  MainCenter,
			AlignCross: CrossCenter,
			Children: []base.Widget{
				&Label{Text: text},
				&Slider{
					Value: value,
					OnChange: func(v float64) {
						value = v
						update()
					},
				},
			},
		}
	}

	err := loop.Run(func() error {
		w, err := windows.NewWindow("Slider", render())
		if err != nil {
			return err
		}

		mainWindow = w
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("OK")
	}
}

func TestSliderMount(t *testing.T) {
	testMountWidgets(t,
		&Slider{Value: 50},
		&Slider{Value: 10},
		&Slider{Value: 0},
		&Slider{Value: 100},
		&Slider{Value: 50, Disabled: true},
		&Slider{Value: 500, Max: 1000},
		&Slider{Value: 20},
		&Slider{Value: 200, Max: 300},
	)
}

func TestSliderClose(t *testing.T) {
	testCloseWidgets(t,
		&Slider{Value: 50},
		&Slider{Value: 50, Disabled: true},
		&Slider{Value: 500, Max: 1000},
	)
}

func TestSliderFocus(t *testing.T) {
	testCheckFocusAndBlur(t,
		&Slider{Value: 50},
		&Slider{Value: 40},
		&Slider{Value: 500, Max: 1000},
	)
}

func TestSliderUpdate(t *testing.T) {
	testUpdateWidgets(t, []base.Widget{
		&Slider{Value: 50},
		&Slider{Value: 50, Disabled: true},
		&Slider{Value: 500, Max: 1000},
	}, []base.Widget{
		&Slider{Value: 60},
		&Slider{Value: 50, Min: 10, Max: 60},
		&Slider{Value: 500, Max: 1000, Disabled: true},
	})
}

func TestSlider_UpdateValue(t *testing.T) {
	cases := []struct {
		value    float64
		min, max float64
		out      float64
	}{
		{1, 0, 10, 1},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
		{-1, 0, 10, 0},
		{11, 0, 10, 10},
		{-1, 0, 0, 0},
		{11, 0, 0, 0},
		{-1, 0, -1, 0},
	}

	for i, v := range cases {
		slider := Slider{Value: v.value, Min: v.min, Max: v.max}
		slider.UpdateValue()
		if slider.Value != v.out {
			t.Errorf("Case %d: .Value does not match, got %f, want %f", i, slider.Value, v.out)
		}
	}
}

func TestSliderLayout(t *testing.T) {
	testLayoutWidget(t, &Slider{Value: 50, Max: 100})
}

func TestSliderMinSize(t *testing.T) {
	testMinSizeWidget(t, &Slider{Value: 50, Max: 100})
}
