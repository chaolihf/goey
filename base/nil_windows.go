package base

import (
	"github.com/chaolihf/win"
)

// SetOrder is called to ensure that windows appears in the correct order.
// This method is part of the method set required to implement the Element
// interface on WIN32.
func (*nilElement) SetOrder(previous win.HWND) win.HWND {
	return previous
}
