package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestProgressMount(t *testing.T) {
	testMountWidgets(t,
		&Progress{Value: 50, Min: 0, Max: 100},
		&Progress{Value: 0},
		&Progress{Value: 10, Min: 0, Max: 1000},
		&Progress{Value: 0},
		&Progress{Value: 100},
		&Progress{Value: 500, Max: 1000},
		&Progress{Min: 50, Max: 50}, // Identical Min and Max
	)
}

func TestProgressClose(t *testing.T) {
	testCloseWidgets(t,
		&Progress{Value: 50, Min: 0, Max: 100},
		&Progress{Value: 0},
		&Progress{Value: 10, Min: 0, Max: 1000},
	)
}

func TestProgressUpdate(t *testing.T) {
	testUpdateWidgets(t, []base.Widget{
		&Progress{Value: 50, Min: 0, Max: 100},
		&Progress{Value: 50, Min: 0, Max: 100},
		&Progress{Value: 50, Min: 50, Max: 50},
	}, []base.Widget{
		&Progress{Value: 75, Min: 0, Max: 100},
		&Progress{Value: 50, Min: 0, Max: 200},
		&Progress{Value: 150, Min: 100, Max: 200},
	})
}

func TestProgress_UpdateValue(t *testing.T) {
	cases := []struct {
		value    int
		min, max int
		out      int
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
		widget := Progress{Value: v.value, Min: v.min, Max: v.max}
		widget.UpdateValue()
		if widget.Value != v.out {
			t.Errorf("Case %d: .Value does not match, got %d, want %d", i, widget.Value, v.out)
		}
	}
}
