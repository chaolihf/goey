package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestHRMount(t *testing.T) {
	testingMountWidgets(t,
		&HR{},
		&HR{},
		&HR{},
	)
}

func TestHRClose(t *testing.T) {
	testingCloseWidgets(t,
		&HR{},
		&HR{},
	)
}

func TestHRUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&HR{},
		&HR{},
	}, []base.Widget{
		&HR{},
		&HR{},
	})
}

func TestHRLayout(t *testing.T) {
	cases := []struct {
		name string
		bc   base.Constraints
	}{
		{"expand", base.Expand()},
		{"expand-height", base.ExpandHeight(96 * DIP)},
		{"expand-width", base.ExpandWidth(24 * DIP)},
		{"loose", base.Loose(base.Size{96 * DIP, 24 * DIP})},
		{"tight", base.Tight(base.Size{96 * DIP, 24 * DIP})},
		{"tight-height", base.TightHeight(24 * DIP)},
		{"tight-width", base.TightWidth(96 * DIP)},
	}

	updater, closer := testingLayoutWidget(t, &HR{})
	defer closer()

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			size := updater(v.bc)
			if !v.bc.IsSatisfiedBy(size) {
				t.Errorf("layout does not respect constraints")
			}
		})
	}
}

func TestHRMinSize(t *testing.T) {
	testingMinSizeWidget(t, &HR{})
}
