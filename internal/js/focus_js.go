package goeyjs

import (
	"syscall/js"
)

type FocusCB struct {
	callback
	Fn func()
}

func (cb *FocusCB) Set(elem js.Value, onfocus func()) {
	cb.Fn = onfocus

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			cb.Fn()
			return nil
		})
		elem.Set("onfocus", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Set("onfocus", js.Undefined())
	}
}

type BlurCB struct {
	callback
	Fn func()
}

func (cb *BlurCB) Set(elem js.Value, onblur func()) {
	cb.Fn = onblur

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			cb.Fn()
			return nil
		})
		elem.Set("onblur", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Set("onblur", js.Undefined())
	}
}
