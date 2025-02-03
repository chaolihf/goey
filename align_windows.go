package goey

import (
	"github.com/chaolihf/win"
)

func (w *AlignElement) SetOrder(previous win.HWND) win.HWND {
	previous = w.child.SetOrder(previous)
	return previous
}
