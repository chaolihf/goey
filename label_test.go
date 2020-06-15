package goey

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"bitbucket.org/rj/goey/base"
)

func labelValues(values []reflect.Value, rand *rand.Rand) {
	const complexSize = 50

	// This is copied from the testing/quick package, but modified somewhat.
	// The function in the standard library will create strings using all
	// code points in the range up to 0x10FFFF.  This works fine on Linux,
	// but on Windows unrecognized codepoints are replaced with 0xFFFD,
	// which is appropriate but breaks the tests.  Here, we restrict code
	// points to ASCII less the control characters.
	numChars := rand.Intn(complexSize)
	codePoints := make([]rune, numChars)
	for i := 0; i < numChars; i++ {
		codePoints[i] = rune(0x20 + rand.Intn(0x7F-0x20))
	}
	values[0] = reflect.ValueOf(string(codePoints))
}

func TestLabelMount(t *testing.T) {
	testingMountWidgets(t,
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
		&Label{Text: ""},
		&Label{Text: "ABCD\nEDFG"},
	)

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		f := func(text string) bool {
			return testingMountWidget(t, &Label{Text: text})
		}
		if err := quick.Check(f, &quick.Config{Values: labelValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}

func TestLabelClose(t *testing.T) {
	testingCloseWidgets(t,
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
	)
}

func TestLabelUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
		&Label{Text: ""},
		&Label{Text: "ABCD\nEDFG"},
	}, []base.Widget{
		&Label{Text: ""},
		&Label{Text: "ABCD\nEDFG"},
		&Label{Text: "AB"},
		&Label{Text: "BC"},
		&Label{Text: "CD"},
	})

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		updater, closer := testingUpdateWidget(t)
		defer closer()

		f := func(text string) bool {
			return updater(&Label{Text: text})
		}
		if err := quick.Check(f, &quick.Config{Values: labelValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}

func TestLabelLayout(t *testing.T) {
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

	updater, closer := testingLayoutWidget(t, &Label{Text: "AB"})
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

func TestLabelMinSize(t *testing.T) {
	testingMinSizeWidget(t, &Label{Text: "AB"})
}
