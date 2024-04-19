//go:build go1.12
// +build go1.12

package goey

import (
	"strconv"
	"syscall/js"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/internal/js"
)

type tabsElement struct {
	Control
	innerDiv js.Value
	clickCB  js.Func

	value    int
	child    base.Element
	widgets  []TabItem
	insets   Insets
	onChange func(int)

	cachedInsets base.Point
	cachedBounds base.Rectangle
	cachedTabsW  base.Length
}

func (w *Tabs) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := goeyjs.CreateElement("ul", "goey nav nav-tabs")
	innerDiv := goeyjs.CreateElement("div", "goey goey-tabs-panel")
	parent.Handle.Call("appendChild", handle)
	parent.Handle.Call("appendChild", innerDiv)

	// Create the element
	retval := &tabsElement{
		Control:  Control{handle},
		innerDiv: innerDiv,
		insets:   w.Insets,

		value:   len(w.Children), // Force tab change
		widgets: w.Children,
	}
	retval.attachOnClick()
	retval.updateProps(w)

	return retval, nil
}

func (w *tabsElement) attachOnClick() {
	w.clickCB = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		value, _ := strconv.Atoi(args[0].Get("currentTarget").Get("dataset").Get("value").String())
		w.onClick(value)
		return nil
	})
}

func (w *tabsElement) contentInsets() base.Point {
	// Padding for the tab panel should match the padding used by Bootstrap for
	// the tabs.  Also add 1px for the borders.

	if w.cachedInsets.Y == 0 {
		w.cachedInsets = base.Point{
			X: base.FromPixelsX(1 + 1),
			Y: base.FromPixelsY(40 + 1),
		}
	}

	return w.cachedInsets
}

func (w *tabsElement) controlTabsMinWidth() base.Length {
	if w.cachedTabsW == 0 {
		// TODO:
		w.cachedTabsW = base.FromPixelsX(32)
	}
	return w.cachedTabsW
}

func (w *tabsElement) mountPage(page int) error {
	child, err := base.Mount(base.Control{w.innerDiv}, w.widgets[page].Child)
	if err != nil {
		return err
	}
	if w.cachedBounds.Dx() != 0 {
		child.Layout(base.Tight(base.Size{
			Width:  w.cachedBounds.Dx(),
			Height: w.cachedBounds.Dy(),
		}))

		child.SetBounds(w.cachedBounds)
	}

	if w.child != nil {
		w.child.Close()
	}
	w.child = child

	// Update the tabs
	if w.value < len(w.widgets) {
		a := w.handle.Get("children").Index(w.value).Get("children").Index(0)
		a.Set("className", "nav-link")
		a.Set("aria-current", js.Undefined())
	}
	{
		a := w.handle.Get("children").Index(page).Get("children").Index(0)
		a.Set("className", "nav-link active")
		a.Set("aria-current", "page")
	}

	return nil
}

func (w *tabsElement) onClick(value int) {
	if value != w.value {
		if w.onChange != nil {
			w.onChange(value)
		}
		if value != w.value {
			// Not clear how an error at this point should be handled.
			// The widget is supposed to already be mounted, but we create and
			// remove controls when the tab is changed.
			// In practice, errors are very infrequent (never?).  GTK widgets
			// will never fail to mount.
			_ = w.mountPage(value)
			w.value = value
		}
	}
}

func (w *tabsElement) Props() base.Widget {
	count := w.handle.Get("childElementCount").Int()

	children := make([]TabItem, count)
	for i := 0; i < count; i++ {
		children[i].Caption = w.handle.Get("children").Index(i).Get("children").Index(0).Get("textContent").String()
		children[i].Child = w.widgets[i].Child
	}

	return &Tabs{
		Value:    w.value,
		Children: children,
		Insets:   w.insets,
		OnChange: w.onChange,
	}
}

func (w *tabsElement) SetBounds(bounds base.Rectangle) {
	calcTabStrip := func(b base.Rectangle) base.Rectangle {
		b.Max.Y = b.Min.Y + base.FromPixelsY(40)
		return b
	}
	w.Control.SetBounds(calcTabStrip(bounds))

	calcTabPanel := func(b base.Rectangle) base.Rectangle {
		b.Min.Y += base.FromPixelsY(40)
		return b
	}
	(&Control{w.innerDiv}).SetBounds(calcTabPanel(bounds))

	insets := w.contentInsets()
	insets.X += w.insets.Dx()
	insets.Y += w.insets.Dy()
	if bounds.Dx() > insets.X && bounds.Dy() > insets.Y {
		bounds = base.Rectangle{
			Max: base.Point{bounds.Dx() - insets.X, bounds.Dy() - insets.Y},
		}

		// Offset
		offset := base.Point{w.insets.Left, w.insets.Top}
		bounds.Min = bounds.Min.Add(offset)
		bounds.Max = bounds.Max.Add(offset)
	}

	// Update bounds for the child
	// Cache and update child element
	w.cachedBounds = bounds
	w.child.SetBounds(w.cachedBounds)
}

func updateTabItems(handle js.Value, clickCB js.Value, items []TabItem) {
	n := handle.Get("childElementCount").Int()

	// Remove excess options from the element
	children := handle.Get("children")
	if n > len(items) {
		for i := n; i > len(items); i-- {
			children.Index(i - 1).Call("remove")
		}
		n = len(items)
	}

	// Change text of existing options
	for i := 0; i < n; i++ {
		children.Index(i).Get("children").Index(0).Set("textContent", items[i].Caption)
	}

	// Add new options
	for i := n; i < len(items); i++ {
		li := goeyjs.CreateElement("li", "nav-item")
		li.Get("dataset").Set("value", i)
		li.Set("onclick", clickCB)
		a := goeyjs.CreateElement("a", "nav-link")
		a.Set("textContent", items[i].Caption)
		a.Set("href", "#")
		li.Call("appendChild", a)
		handle.Call("appendChild", li)
	}
}

func (w *tabsElement) updateProps(data *Tabs) error {
	updateTabItems(w.handle, w.clickCB.Value, data.Children)
	w.widgets = data.Children

	if w.value != data.Value {
		w.mountPage(data.Value)
		w.value = data.Value
	}

	return nil
}
