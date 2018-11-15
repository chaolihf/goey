// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type textareaElement struct {
	control  *cocoa.TextField
	minLines int
}

func (w *TextArea) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewTextField(parent.Handle, w.Value)
	control.SetValue(w.Value)
	control.SetPlaceholder(w.Placeholder)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur)

	retval := &textareaElement{
		control:  control,
		minLines: w.MinLines,
	}
	return retval, nil
}

func (w *textareaElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *textareaElement) Layout(bc base.Constraints) base.Size {
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *textareaElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *textareaElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *textareaElement) Props() base.Widget {
	onchange, onfocus, onblur := w.control.Callbacks()

	return &TextArea{
		Value:       w.control.Value(),
		Disabled:    !w.control.IsEnabled(),
		Placeholder: w.control.Placeholder(),
		ReadOnly:    false,
		MinLines:    w.minLines,
		OnChange:    onchange,
		OnFocus:     onfocus,
		OnBlur:      onblur,
	}
}

func (w *textareaElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *textareaElement) TakeFocus() bool {
	return w.control.MakeFirstResponder()
}

func (w *textareaElement) updateProps(data *TextArea) error {
	w.control.SetValue(data.Value)
	w.control.SetPlaceholder(data.Placeholder)
	w.control.SetEnabled(!data.Disabled)
	w.minLines = data.MinLines
	w.control.SetCallbacks(data.OnChange, data.OnFocus, data.OnBlur)
	return nil
}
