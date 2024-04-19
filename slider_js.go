//go:build go1.12
// +build go1.12

package goey

import (
	"strconv"

	"github.com/chaolihf/goey/base"
	"github.com/chaolihf/goey/internal/js"
)

type sliderElement struct {
	Control

	onChange goeyjs.ChangeFloat64CB
	onFocus  goeyjs.FocusCB
	onBlur   goeyjs.BlurCB
}

func (w *Slider) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := goeyjs.CreateElement("input", "goey form-range")
	handle.Set("type", "range")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &sliderElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *sliderElement) Props() base.Widget {
	value, _ := strconv.ParseFloat(w.handle.Get("value").String(), 64)
	min, _ := strconv.ParseFloat(w.handle.Get("min").String(), 64)
	max, _ := strconv.ParseFloat(w.handle.Get("max").String(), 64)

	return &Slider{
		Value:    value,
		Disabled: w.handle.Get("disabled").Truthy(),
		Min:      min,
		Max:      max,
		OnChange: w.onChange.Fn,
		OnFocus:  w.onFocus.Fn,
		OnBlur:   w.onBlur.Fn,
	}
}

func (w *sliderElement) updateProps(data *Slider) error {
	w.handle.Set("min", data.Min)
	w.handle.Set("max", data.Max)
	w.handle.Set("value", data.Value)
	w.handle.Set("disabled", data.Disabled)
	w.onChange.Set(w.handle, data.OnChange)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
