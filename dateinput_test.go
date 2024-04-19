package goey

import (
	"testing"
	"time"

	"github.com/chaolihf/goey/base"
)

func TestDateInputMount(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 0, 0, 0, 0, time.Local)
	v2 := time.Date(2007, time.February, 3, 0, 0, 0, 0, time.Local)
	v3 := time.Date(2007, time.March, 4, 0, 0, 0, 0, time.Local)

	testMountWidgets(t,
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v3},
	)
}

func TestDateInputClose(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 0, 0, 0, 0, time.Local)
	v2 := time.Date(2007, time.February, 3, 0, 0, 0, 0, time.Local)
	v3 := time.Date(2007, time.March, 4, 0, 0, 0, 0, time.Local)

	testCloseWidgets(t,
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v3},
	)
}

func TestDateInputEvents(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 15, 4, 5, 0, time.Local)

	testCheckFocusAndBlur(t,
		&DateInput{Value: v1},
		&DateInput{Value: v2},
		&DateInput{Value: v2},
	)
}

func TestDateInputUpdate(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 0, 0, 0, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 0, 0, 0, 0, time.Local)

	testUpdateWidgets(t, []base.Widget{
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v2},
	}, []base.Widget{
		&DateInput{Value: v2},
		&DateInput{Value: v2, Disabled: false},
		&DateInput{Value: v1, Disabled: true},
	})
}
