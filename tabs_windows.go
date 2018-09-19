package goey

import (
	"syscall"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/lxn/win"
)

var (
	tabs struct {
		className     []uint16
		oldWindowProc uintptr
	}
)

func init() {
	tabs.className = []uint16{'S', 'y', 's', 'T', 'a', 'b', 'C', 'o', 'n', 't', 'r', 'o', 'l', '3', '2', 0}
}

func (w *Tabs) mount(parent base.Control) (base.Element, error) {
	style := uint32(win.WS_CHILD | win.WS_VISIBLE)
	hwnd := win.CreateWindowEx(0, &tabs.className[0], nil, style,
		10, 10, 100, 100,
		parent.HWnd, win.HMENU(nextControlID()), 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hMessageFont), 0)
	}

	for i, v := range w.Children {
		text, err := syscall.UTF16PtrFromString(v.Caption)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		item := win.TCITEM{
			Mask:    win.TCIF_TEXT,
			PszText: text,
			LParam:  uintptr(i),
		}
		win.SendMessage(hwnd, win.TCM_INSERTITEM, uintptr(i), uintptr(unsafe.Pointer(&item)))
	}
	if w.Value > 0 {
		win.SendMessage(hwnd, win.TCM_SETCURSEL, uintptr(w.Value), 0)
	}

	child := base.Element(nil)
	if len(w.Children) > 0 {
		err := error(nil)
		if w.Value >= 0 {
			child, err = base.Mount(parent, w.Children[w.Value].Child)
		} else {
			child, err = base.Mount(parent, w.Children[0].Child)
		}
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}
	}

	retval := &tabsElement{
		Control:  Control{hwnd},
		child:    child,
		parent:   parent,
		value:    w.Value,
		insets:   w.Insets,
		widgets:  w.Children,
		onChange: w.OnChange,
	}

	// Subclass the window procedure
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))
	subclassWindowProcedure(hwnd, &tabs.oldWindowProc, syscall.NewCallback(tabsWindowProc))

	return retval, nil
}

type tabsElement struct {
	Control
	child    base.Element
	parent   base.Control
	value    int
	insets   Insets
	widgets  []TabItem
	onChange func(int)

	cachedInsets base.Point
	cachedBounds base.Rectangle
}

func (w *tabsElement) controlInsets() base.Point {
	if w.cachedInsets.Y == 0 {
		rect := win.RECT{}

		win.SendMessage(w.hWnd, win.TCM_ADJUSTRECT, win.TRUE, uintptr(unsafe.Pointer(&rect)))
		w.cachedInsets = base.Point{
			X: base.FromPixelsX(int(rect.Right - rect.Left)),
			Y: base.FromPixelsY(int(rect.Bottom - rect.Top)),
		}
	}

	return w.cachedInsets
}

func (w *tabsElement) Props() base.Widget {
	count := win.SendMessage(w.hWnd, win.TCM_GETITEMCOUNT, 0, 0)
	children := make([]TabItem, count)
	for i := uintptr(0); i < count; i++ {
		text := [128]uint16{}
		item := win.TCITEM{
			Mask:       win.TCIF_TEXT,
			PszText:    &text[0],
			CchTextMax: 128,
		}

		win.SendMessage(w.hWnd, win.TCM_GETITEM, i, uintptr(unsafe.Pointer(&item)))
		children[i].Caption = syscall.UTF16ToString(text[:])
		children[i].Child = w.widgets[i].Child
	}

	return &Tabs{
		Value:    int(win.SendMessage(w.hWnd, win.TCM_GETCURSEL, 0, 0)),
		Children: children,
		OnChange: w.onChange,
	}
}

func (w *tabsElement) SetOrder(previous win.HWND) win.HWND {
	previous = w.Control.SetOrder(previous)
	if w.child != nil {
		previous = w.child.SetOrder(previous)
	}
	return previous
}

func (w *tabsElement) SetBounds(bounds base.Rectangle) {
	w.Control.SetBounds(bounds)

	if w.child != nil {
		// Determine the bounds for the child widget
		rect := win.RECT{}
		win.SendMessage(w.hWnd, win.TCM_ADJUSTRECT, win.FALSE, uintptr(unsafe.Pointer(&rect)))
		w.cachedBounds = base.Rectangle{
			Min: bounds.Min.Add(base.Point{base.FromPixelsX(int(rect.Left)), base.FromPixelsY(int(rect.Top))}),
			Max: bounds.Max.Add(base.Point{base.FromPixelsX(int(rect.Right)), base.FromPixelsY(int(rect.Bottom))}),
		}
		// Offset to handle insets
		w.cachedBounds.Min.X += w.insets.Left
		w.cachedBounds.Min.Y += w.insets.Top
		w.cachedBounds.Max.X -= w.insets.Right
		w.cachedBounds.Max.Y -= w.insets.Bottom

		// Update bounds for the child
		w.child.SetBounds(w.cachedBounds)
	}
}

func (w *tabsElement) updateChildren(children []TabItem) error {
	len1 := len(w.widgets)
	len2 := len(children)

	if len1 <= len2 {
		// Change caption for tabs that already exist
		for i, v := range children[:len1] {
			text, err := syscall.UTF16PtrFromString(v.Caption)
			if err != nil {
				return err
			}

			item := win.TCITEM{
				Mask:    win.TCIF_TEXT,
				PszText: text,
			}
			win.SendMessage(w.hWnd, win.TCM_SETITEM, uintptr(i), uintptr(unsafe.Pointer(&item)))
		}

		// Add new tabs to extend the list
		for i, v := range children[len1:] {
			text, err := syscall.UTF16PtrFromString(v.Caption)
			if err != nil {
				return err
			}

			item := win.TCITEM{
				Mask:    win.TCIF_TEXT,
				PszText: text,
			}
			win.SendMessage(w.hWnd, win.TCM_INSERTITEM, uintptr(i+len1), uintptr(unsafe.Pointer(&item)))
		}
	} else {
		// Change caption for tabs that already exist
		for i, v := range children {
			text, err := syscall.UTF16PtrFromString(v.Caption)
			if err != nil {
				return err
			}

			item := win.TCITEM{
				Mask:    win.TCIF_TEXT,
				PszText: text,
			}
			win.SendMessage(w.hWnd, win.TCM_SETITEM, uintptr(i), uintptr(unsafe.Pointer(&item)))
		}

		// Delete excess tabs.
		for i := len2; i < len1; i++ {
			win.SendMessage(w.hWnd, win.TCM_DELETEITEM, uintptr(i), 0)
		}
	}

	w.widgets = children
	return nil
}

func (w *tabsElement) updateProps(data *Tabs) error {
	// Update the tabs
	err := w.updateChildren(data.Children)
	if err != nil {
		return err
	}

	// Update which tab is currently selected
	if data.Value >= 0 && w.value != data.Value {
		win.SendMessage(w.hWnd, win.TCM_SETCURSEL, uintptr(data.Value), 0)
		w.value = data.Value
	}

	// Update event handlers
	w.onChange = data.OnChange

	return nil
}

func tabsWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		tabsGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_NOTIFY:
		if n := (*win.NMHDR)(unsafe.Pointer(lParam)); true {
			if n.Code == uint32(0x100000000+win.TCN_SELCHANGE) {
				cursel := int(win.SendMessage(hwnd, win.TCM_GETCURSEL, 0, 0))
				if w := tabsGetPtr(hwnd); w.value != cursel {
					if w.onChange != nil {
						w.onChange(cursel)
					}
					if w.value != cursel {
						child, err := base.DiffChild(w.parent, w.child, w.widgets[cursel].Child)
						if err != nil {
							panic("Unhandled error!")
						}
						if child != nil {
							child.SetOrder(w.hWnd)
							child.Layout(base.Tight(base.Size{
								Width:  w.cachedBounds.Dx(),
								Height: w.cachedBounds.Dy(),
							}))
							child.SetBounds(w.cachedBounds)
							win.InvalidateRect(win.GetParent(w.hWnd), nil, false)
						}
						w.child = child
						w.value = cursel
					}
				}
			}
		}
		return 0
	}

	return win.CallWindowProc(tabs.oldWindowProc, hwnd, msg, wParam, lParam)
}

func tabsGetPtr(hwnd win.HWND) *tabsElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*tabsElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
