package goey

import (
	"github.com/chaolihf/win"
)

func (w *emptyElement) SetOrder(hwnd win.HWND) win.HWND {
	return hwnd
}
