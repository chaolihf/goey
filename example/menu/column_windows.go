package main

import (
	"github.com/chaolihf/win"
)

func (w *columnElement) SetOrder(previous win.HWND) win.HWND {
	for _, v := range w.children {
		previous = v.SetOrder(previous)
	}
	return previous
}
