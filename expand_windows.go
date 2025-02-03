package goey

import (
	"github.com/chaolihf/win"
)

func (w *ExpandElement) SetOrder(previous win.HWND) win.HWND {
	if w.child != nil {
		previous = w.child.SetOrder(previous)
	}
	return previous
}
