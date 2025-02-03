package animate

import (
	"github.com/chaolihf/win"
)

func (w *wipeElement) SetOrder(previous win.HWND) win.HWND {
	return w.child.SetOrder(previous)
}

func (w *wipeElement) paint() {
	win.InvalidateRect(w.parent.HWnd, nil, true)

}
