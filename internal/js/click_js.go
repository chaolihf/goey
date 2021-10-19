package goeyjs

import (
	"syscall/js"

	"gitlab.com/stone.code/assert"
)

type ClickCB struct {
	callback
	Fn func()
}

func (cb *ClickCB) Set(elem js.Value, onclick func()) {
	assert.Assert((cb.Fn != nil) == cb.jsfunc.Truthy(), "callback not syncrhonized")

	cb.Fn = onclick

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			cb.Fn()
			return nil
		})
		elem.Set("onclick", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Set("onclick", js.Undefined())
	}
}
