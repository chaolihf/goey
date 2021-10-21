package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestSelectInputMount(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testMountWidgets(t,
		&SelectInput{Value: 0, Items: options},
		&SelectInput{Value: 1, Items: options},
		&SelectInput{Value: 2, Items: options, Disabled: true},
		&SelectInput{Unset: true, Items: options, Disabled: true},
		&SelectInput{Items: []string{}},
	)
}

func TestSelectInputClose(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testCloseWidgets(t,
		&SelectInput{Value: 0, Items: options},
		&SelectInput{Value: 1, Items: options},
		&SelectInput{Value: 2, Items: options, Disabled: true},
		&SelectInput{Unset: true, Items: options, Disabled: true},
	)
}

func TestSelectInputEvents(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testCheckFocusAndBlur(t,
		&SelectInput{Items: options},
		&SelectInput{Items: options},
		&SelectInput{Items: options},
	)
}

func TestSelectInputUpdate(t *testing.T) {
	options1 := []string{"Option A", "Option B", "Option C"}
	options2 := []string{"Choice A", "Choice B", "Choice C"}

	testUpdateWidgets(t, []base.Widget{
		&SelectInput{Value: 0, Items: options1},
		&SelectInput{Value: 1, Items: options2},
		&SelectInput{Value: 2, Items: options1, Disabled: true},
		&SelectInput{Unset: true, Items: options2},
	}, []base.Widget{
		&SelectInput{Value: 1, Items: options2},
		&SelectInput{Unset: true, Items: options1},
		&SelectInput{Value: 2, Items: options1, Disabled: true},
		&SelectInput{Value: 1, Items: options2},
	})
}

func TestSelectInputLayout(t *testing.T) {
	testLayoutWidget(t, &SelectInput{
		Items: []string{"Option A", "Option B", "Option C"},
	})
}

func TestSelectInputMinSize(t *testing.T) {
	testMinSizeWidget(t, &SelectInput{
		Items: []string{"Option A", "Option B", "Option C"},
	})
}

func TestSelectInput_UpdateValue(t *testing.T) {
	cases := []struct {
		value int
		out   int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 2},
	}

	for i, v := range cases {
		input := &SelectInput{
			Value: v.value,
			Items: []string{"Option A", "Option B", "Option C"},
		}
		input.UpdateValue()
		if input.Value != v.out {
			t.Errorf("Case %d: .Value does not match, got %d, want %d", i, input.Value, v.out)
		}
	}
}
