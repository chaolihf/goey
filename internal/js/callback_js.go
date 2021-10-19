package goeyjs

import "syscall/js"

type callback struct {
	jsfunc js.Func
}

func (cb *callback) release() {
	cb.jsfunc.Release()
	cb.jsfunc.Value = js.Undefined()
}

func (cb *callback) Close() {
	cb.jsfunc.Release()
}
