package goeyjs

import (
	"strconv"
	"syscall/js"
	"time"

	"gitlab.com/stone.code/assert"
)

type ChangeStringCB struct {
	callback
	Fn func(string)
}

func (cb *ChangeStringCB) Set(elem js.Value, onchange func(string)) {
	assert.Assert((cb.Fn != nil) == cb.jsfunc.Truthy(), "callback not syncrhonized")

	cb.Fn = onchange

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			cb.Fn(elem.Get("value").String())
			return nil
		})
		elem.Set("oninput", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Delete("oninput")
	}
}

type ChangeBoolCB struct {
	callback
	Fn func(bool)
}

func (cb *ChangeBoolCB) Set(elem js.Value, onchange func(bool)) {
	assert.Assert((cb.Fn != nil) == cb.jsfunc.Truthy(), "callback not syncrhonized")

	cb.Fn = onchange

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			cb.Fn(elem.Get("checked").Truthy())
			return nil
		})
		elem.Set("oninput", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Delete("oninput")
	}
}

type ChangeIntCB struct {
	callback
	Fn func(int)
}

func (cb *ChangeIntCB) Set(elem js.Value, onchange func(int)) {
	assert.Assert((cb.Fn != nil) == cb.jsfunc.Truthy(), "callback not syncrhonized")

	cb.Fn = onchange

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			value, err := strconv.Atoi(elem.Get("value").String())
			assert.Assert(err == nil, "value of HTMLSelectInput did not convert to int")
			cb.Fn(value)
			return nil
		})
		elem.Set("oninput", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Delete("oninput")
	}
}

type ChangeInt64CB struct {
	callback
	Fn func(int64)
}

func (cb *ChangeInt64CB) Set(elem js.Value, onchange func(int64)) {
	assert.Assert((cb.Fn != nil) == cb.jsfunc.Truthy(), "callback not syncrhonized")

	cb.Fn = onchange

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			value, err := strconv.ParseInt(elem.Get("value").String(), 10, 64)
			assert.Assert(err == nil, "value of HTMLSelectInput did not convert to int")
			cb.Fn(value)
			return nil
		})
		elem.Set("oninput", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Delete("oninput")
	}
}

type ChangeDateCB struct {
	callback
	Fn func(time.Time)
}

func (cb *ChangeDateCB) Set(elem js.Value, onchange func(time.Time)) {
	assert.Assert((cb.Fn != nil) == cb.jsfunc.Truthy(), "callback not syncrhonized")

	cb.Fn = onchange

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			s := elem.Get("value").String()
			if s != "" {
				value, err := time.Parse("2006-1-2", s)
				assert.Assert(err == nil, "value of HTMLInput did not convert to date")
				cb.Fn(value)
			}
			return nil
		})
		elem.Set("oninput", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Set("oninput", js.Undefined())
	}
}

type ChangeFloat64CB struct {
	callback
	Fn func(float64)
}

func (cb *ChangeFloat64CB) Set(elem js.Value, onchange func(float64)) {
	assert.Assert((cb.Fn != nil) == cb.jsfunc.Truthy(), "callback not syncrhonized")

	cb.Fn = onchange

	if cb.Fn != nil && cb.jsfunc.IsUndefined() {
		cb.jsfunc = js.FuncOf(func(js.Value, []js.Value) interface{} {
			value, err := strconv.ParseFloat(elem.Get("value").String(), 64)
			assert.Assert(err == nil, "value of HTMLSelectInput did not convert to int")
			cb.Fn(value)
			return nil
		})
		elem.Set("oninput", cb.jsfunc)
	} else if cb.Fn == nil && !cb.jsfunc.IsUndefined() {
		cb.release()
		elem.Set("oninput", js.Undefined())
	}
}
