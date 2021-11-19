package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestTabsMount(t *testing.T) {
	items := []TabItem{
		{"Tab 1", &Button{Text: "Click me A!"}},
		{"Tab 2", &Button{Text: "Click me B!"}},
	}
	emptyTabs := []TabItem{
		{"Tab A", nil},
		{"Tab B", nil},
		{"Tab C", nil},
	}

	testMountWidgets(t,
		&Tabs{Children: items},
		&Tabs{Value: 1, Children: items},
		&Tabs{Value: 2, Children: emptyTabs},
	)
}

func TestTabsClose(t *testing.T) {
	items := []TabItem{
		{"Tab 1", &Button{Text: "Click me A!"}},
		{"Tab 2", &Button{Text: "Click me B!"}},
	}

	testCloseWidgets(t,
		&Tabs{Children: items},
		&Tabs{Value: 1, Children: items},
	)
}

func TestTabsUpdate(t *testing.T) {
	items1 := []TabItem{
		{"Tab 1", &Button{Text: "Click me!"}},
		{"Tab 2", &Button{Text: "Don't click me!"}},
	}
	items2 := []TabItem{
		{"Tab A", &Button{Text: "Don't click me!"}},
		{"Tab B", &Button{Text: "Click me!"}},
		{"Tab C", &Button{Text: "..."}},
	}

	testUpdateWidgets(t, []base.Widget{
		&Tabs{Children: items1},
		&Tabs{Value: 1, Children: items2},
	}, []base.Widget{
		&Tabs{Value: 1, Children: items2},
		&Tabs{Children: items1},
	})
}

func TestTabsLayout(t *testing.T) {
	testLayoutWidget(t, &Tabs{Children: []TabItem{
		{"Tab A", &Button{Text: "Don't click me!"}},
		{"Tab B", &Button{Text: "Click me!"}},
		{"Tab C", &Button{Text: "..."}},
	}})
}

func TestTabsMinSize(t *testing.T) {
	testMinSizeWidget(t, &Tabs{Children: []TabItem{
		{"Tab A", &Button{Text: "Don't click me!"}},
		{"Tab B", &Button{Text: "Click me!"}},
		{"Tab C", &Button{Text: "..."}},
	}})
}

func TestTabs_UpdateValue(t *testing.T) {
	someTabs := []TabItem{
		{"Tab A", nil},
		{"Tab B", nil},
		{"Tab C", nil},
	}

	cases := []struct {
		value int
		tabs  []TabItem
		out   int
	}{
		{0, someTabs, 0},
		{2, someTabs, 2},
		{-10, someTabs, 0},
		{10, someTabs, 2},
		{0, nil, -1},
	}

	for i, v := range cases {
		widget := &Tabs{
			Value:    v.value,
			Children: v.tabs,
		}

		widget.UpdateValue()
		if widget.Value != v.out {
			t.Errorf("case %d: want %d, got %d", i, v.out, widget.Value)
		}
	}
}
