package goeyjs

import (
	"strconv"
	"syscall/js"

	"gitlab.com/stone.code/assert"
)

type EnterKeyCB struct {
	callback
	Fn func(string)
}

func (cb *EnterKeyCB) Set(elem js.Value, onenterkey func(string)) {
	cb.Fn = onenterkey

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			event := args[0]
			if event.Get("keyCode").Int() == 13 {
				event.Call("preventDefault")
				cb.Fn(elem.Get("value").String())
			}
			return nil
		})
		elem.Set("onkeyup", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Set("onkeyup", js.Undefined())
	}
}

type EnterKeyInt64CB struct {
	callback
	Fn func(int64)
}

func (cb *EnterKeyInt64CB) Set(elem js.Value, onenterkey func(int64)) {
	cb.Fn = onenterkey

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			event := args[0]
			if event.Get("keyCode").Int() == 13 {
				event.Call("preventDefault")
				value, err := strconv.ParseInt(elem.Get("value").String(), 10, 64)
				assert.Assert(err == nil, "failed to convert string to int")
				cb.Fn(value)
			}
			return nil
		})
		elem.Set("onkeyup", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Set("onkeyup", js.Undefined())
	}
}
