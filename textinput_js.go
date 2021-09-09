package goey

import (
	"syscall/js"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/js"
)

type textinputElement struct {
	Control

	onChange   goeyjs.ChangeStringCB
	onFocus    goeyjs.FocusCB
	onBlur     goeyjs.BlurCB
	onEnterKey goeyjs.EnterKeyCB
}

func (w *TextInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "input")
	handle.Get("style").Set("position", "absolute")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &textinputElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *textinputElement) Props() base.Widget {
	return &TextInput{
		Value:       w.handle.Get("value").String(),
		Placeholder: w.handle.Get("placeholder").String(),
		Disabled:    w.handle.Get("disabled").Truthy(),
		Password:    w.handle.Get("type").String() == "password",
		ReadOnly:    w.handle.Get("readonly").Truthy(),
		OnChange:    w.onChange.Fn,
		OnFocus:     w.onFocus.Fn,
		OnBlur:      w.onBlur.Fn,
		OnEnterKey:  w.onEnterKey.Fn,
	}
}

func (w *textinputElement) updateProps(data *TextInput) error {
	w.handle.Set("value", data.Value)
	w.handle.Set("placeholder", data.Placeholder)
	w.handle.Set("disabled", data.Disabled)
	if data.Password {
		w.handle.Set("type", "password")
	} else {
		w.handle.Set("type", "text")
	}
	w.handle.Set("readonly", data.ReadOnly)

	w.onChange.Set(w.handle, data.OnChange)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)
	w.onEnterKey.Set(w.handle, data.OnEnterKey)

	return nil
}
