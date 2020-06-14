package goey

import (
	"reflect"
	"runtime"
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestIntInputMount(t *testing.T) {
	testingMountWidgets(t,
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
		&IntInput{Value: 4, Min: 0, Max: 10},
		&IntInput{Value: 5, Min: -1000, Max: 1000},
	)
}

func TestIntInputClose(t *testing.T) {
	testingCloseWidgets(t,
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
		&IntInput{Value: 4, Min: 0, Max: 10},
		&IntInput{Value: 5, Min: -1000, Max: 1000},
	)
}

func TestIntInputOnFocus(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&IntInput{},
		&IntInput{},
		&IntInput{},
	)
}

func TestIntInputOnChange(t *testing.T) {
	log := make([]int64, 0)

	testingTypeKeys(t, "1234",
		&IntInput{OnChange: func(v int64) {
			log = append(log, v)
		}})

	want := []int64{1, 12, 123, 1234}
	if runtime.GOOS == "linux" {
		// Control does not output events for intermediate typing.
		want = []int64{1234}
	}
	if !reflect.DeepEqual(want, log) {
		t.Errorf("Wanted %v, got %v", want, log)
	}
}

func TestIntInputOnEnterKey(t *testing.T) {
	got := int64(0)

	testingTypeKeys(t, "1234\n",
		&IntInput{OnEnterKey: func(v int64) {
			got = v
		}})

	const want = 1234
	if got != want {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func TestIntInputUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
	}, []base.Widget{
		&IntInput{Value: 1},
		&IntInput{Value: 4, Disabled: true},
		&IntInput{Value: 5, Placeholder: "***"},
	})
}

func TestIntInputUpdateValue(t *testing.T) {
	cases := []struct {
		in, min, max, out int64
	}{
		{0, 10, 20, 10},
		{10, 10, 20, 10},
		{15, 10, 20, 15},
		{20, 10, 20, 20},
		{25, 10, 20, 20},
		{-20, -10, 10, -10},
		{-10, -10, 10, -10},
		{0, -10, 10, 0},
		{10, -10, 10, 10},
		{20, -10, 10, 10},
	}

	for i, v := range cases {
		widget := IntInput{
			Value: v.in,
			Min:   v.min,
			Max:   v.max,
		}

		widget.UpdateValue()
		if got := widget.Value; got != v.out {
			t.Errorf("case %d: got %d, want %d", i, got, v.out)
		}
	}
}
