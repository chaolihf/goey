package goey

import (
	"strconv"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/js"
)

type intinputElement struct {
	Control

	onChange   goeyjs.ChangeInt64CB
	onFocus    goeyjs.FocusCB
	onBlur     goeyjs.BlurCB
	onEnterKey goeyjs.EnterKeyInt64CB
}

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "input")
	handle.Get("style").Set("position", "absolute")
	handle.Set("type", "number")
	handle.Set("step", 1)
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &intinputElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *intinputElement) Props() base.Widget {
	value, _ := strconv.ParseInt(
		w.handle.Get("value").String(),
		10, 64)
	min, _ := strconv.ParseInt(
		w.handle.Get("min").String(),
		10, 64)
	max, _ := strconv.ParseInt(
		w.handle.Get("max").String(),
		10, 64)

	return &IntInput{
		Value:       value,
		Min:         min,
		Max:         max,
		Placeholder: w.handle.Get("placeholder").String(),
		Disabled:    w.handle.Get("disabled").Truthy(),
		OnChange:    w.onChange.Fn,
		OnFocus:     w.onFocus.Fn,
		OnBlur:      w.onBlur.Fn,
		OnEnterKey:  w.onEnterKey.Fn,
	}
}

func (w *intinputElement) updateProps(data *IntInput) error {
	w.handle.Set("value", data.Value)
	w.handle.Set("placeholder", data.Placeholder)
	w.handle.Set("disabled", data.Disabled)
	w.handle.Set("min", data.Min)
	w.handle.Set("max", data.Max)

	w.onChange.Set(w.handle, data.OnChange)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
