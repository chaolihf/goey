package goey

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/chaolihf/goey/base"
)

func paragraphValues(values []reflect.Value, rand *rand.Rand) {
	// Get a string
	labelValues(values, rand)

	// Create a random alignment
	values[1] = reflect.ValueOf(TextAlignment(rand.Uint64() % 4))
}

func TestParagraphMount(t *testing.T) {
	testMountWidgets(t,
		&P{Text: "A", Align: JustifyLeft},
		&P{Text: "B", Align: JustifyRight},
		&P{Text: "C", Align: JustifyCenter},
		&P{Text: "D", Align: JustifyFull},
		&P{Text: "", Align: JustifyLeft},
		&P{Text: "ABCD\nEFGH", Align: JustifyLeft},
	)

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		f := func(text string, align TextAlignment) bool {
			return checkMountWidget(t, &P{Text: text, Align: align})
		}
		if err := quick.Check(f, &quick.Config{Values: paragraphValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}

func TestParagraphClose(t *testing.T) {
	testCloseWidgets(t,
		&P{Text: "A", Align: JustifyLeft},
		&P{Text: "B", Align: JustifyRight},
		&P{Text: "C", Align: JustifyCenter},
		&P{Text: "D", Align: JustifyFull},
	)
}

func TestParagraphUpdate(t *testing.T) {
	testUpdateWidgets(t, []base.Widget{
		&P{Text: "A", Align: JustifyLeft},
		&P{Text: "B", Align: JustifyRight},
		&P{Text: "C", Align: JustifyCenter},
		&P{Text: "D", Align: JustifyFull},
		&P{Text: "", Align: JustifyLeft},
		&P{Text: "ABCD\nEFGH", Align: JustifyLeft},
	}, []base.Widget{
		&P{Text: "", Align: JustifyLeft},
		&P{Text: "ABCD\nEFGH", Align: JustifyLeft},
		&P{Text: "AAA", Align: JustifyRight},
		&P{Text: "BAA", Align: JustifyCenter},
		&P{Text: "CAA", Align: JustifyFull},
		&P{Text: "DAA", Align: JustifyLeft},
	})
}

func TestParagraphLayout(t *testing.T) {
	testLayoutWidget(t, &P{Text: "AB"})
}

func TestParagraphMinSize(t *testing.T) {
	testMinSizeWidget(t, &P{Text: "AB"})
}
