package goeyjs

import (
	"syscall/js"

	"gitlab.com/stone.code/assert"
)

type SelectCB struct {
	callback
	Fn func(int)
}

func (cb *SelectCB) Set(elem js.Value, onselect func(int)) {
	const event = "oninput"

	assert.Assert((cb.Fn != nil) == cb.jsfunc.Truthy(), "callback not synchronized")

	cb.Fn = onselect

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			cb.Fn(elem.Get("selectedIndex").Int())
			return nil
		})
		elem.Set(event, cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Delete(event)
	}
}

