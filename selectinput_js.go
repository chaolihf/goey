// +build go1.12

package goey

import (
	"strconv"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	goeyjs "bitbucket.org/rj/goey/internal/js"
)

type selectinputElement struct {
	Control

	onChange goeyjs.ChangeIntCB
	onFocus  goeyjs.FocusCB
	onBlur   goeyjs.BlurCB
}

func (w *SelectInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	document := js.Global().Get("document")
	handle := document.Call("createElement", "select")
	handle.Set("className", "form-control")
	handle.Get("style").Set("position", "absolute")
	opt := document.Call("createElement", "option")
	opt.Set("text", "XXXXXXXX")
	handle.Call("appendChild", opt)
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &selectinputElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *selectinputElement) Close() {
	w.onChange.Close()
	w.onFocus.Close()
	w.onBlur.Close()

	w.Control.Close()
}

func (w *selectinputElement) createMeasurementElement() js.Value {
	document := js.Global().Get("document")

	handle := document.Call("createElement", "select")
	handle.Set("className", "form-control")
	handle.Get("style").Set("visibility", "hidden")

	body := document.Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", handle)

	return handle
}

func (w *selectinputElement) Layout(bc base.Constraints) base.Size {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := base.FromPixelsX(handle.Get("offsetWidth").Int() + 1)
	width = bc.ConstrainWidth(width)
	height := base.FromPixelsY(handle.Get("offsetHeight").Int() + 1)
	height = bc.ConstrainHeight(height)

	return base.Size{width, height}
}

func (w *selectinputElement) MinIntrinsicHeight(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	height := handle.Get("offsetHeight").Int()

	return base.FromPixelsY(height)
}

func (w *selectinputElement) MinIntrinsicWidth(base.Length) base.Length {
	handle := w.createMeasurementElement()
	defer handle.Call("remove")

	width := handle.Get("offsetWidth").Int()

	return base.FromPixelsX(width + 1)
}

func (w *selectinputElement) Props() base.Widget {
	items := []string{}
	n := w.handle.Get("length").Int()
	for i := 0; i < n; i++ {
		items = append(items,
			w.handle.Index(i).Get("text").String())
	}

	si := w.handle.Get("selectedIndex").Int()

	return &SelectInput{
		Items: items,
		Value: func(si int) int {
			if si < 0 {
				return 0
			} else {
				return si
			}
		}(si),
		Unset:    si < 0,
		Disabled: w.handle.Get("disabled").Truthy(),
		OnChange: w.onChange.Fn,
		OnFocus:  w.onFocus.Fn,
		OnBlur:   w.onBlur.Fn,
	}
}

func updateOptionList(handle js.Value, items []string) {
	n := handle.Get("length").Int()

	// Remove excess options from the element
	if n > len(items) {
		for i := n; i > len(items); i-- {
			handle.Call("remove", i-1)
		}
		n = len(items)
	}

	// Change text of existing options
	for i := 0; i < n; i++ {
		handle.Index(i).Set("text", items[i])
	}

	// Add new options
	for i := n; i < len(items); i++ {
		opt := js.Global().Get("document").Call("createElement", "option")
		opt.Set("text", items[i])
		opt.Set("value", strconv.Itoa(i))
		handle.Call("add", opt)
	}
}

func (w *selectinputElement) updateProps(data *SelectInput) error {
	updateOptionList(w.handle, data.Items)

	if data.Unset {
		w.handle.Set("selectedIndex", -1)
	} else {
		w.handle.Set("selectedIndex", data.Value)
	}
	w.handle.Set("disabled", data.Disabled)
	w.onChange.Set(w.handle, data.OnChange)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
