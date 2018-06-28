package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	comboboxClassName     *uint16
	oldComboboxWindowProc uintptr
)

func init() {
	var err error
	comboboxClassName, err = syscall.UTF16PtrFromString("COMBOBOX")
	if err != nil {
		panic(err)
	}
}

func (w *SelectInput) mount(parent Control) (Element, error) {
	if w.Value >= len(w.Items) {
		w.Value = len(w.Items) - 1
	}
	hwnd := win.CreateWindowEx(win.WS_EX_CLIENTEDGE, comboboxClassName, nil,
		win.WS_CHILD|win.WS_VISIBLE|win.WS_TABSTOP|win.CBS_DROPDOWNLIST,
		10, 10, 100, 100,
		parent.hWnd, win.HMENU(nextControlID()), 0, nil)
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

	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Add items to the control
	longestString := ""
	for _, v := range w.Items {
		text, err := syscall.UTF16PtrFromString(v)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}
		win.SendMessage(hwnd, win.CB_ADDSTRING, 0, uintptr(unsafe.Pointer(text)))

		if len(v) > len(longestString) {
			longestString = v
		}
	}
	if !w.Unset {
		win.SendMessage(hwnd, win.CB_SETCURSEL, uintptr(w.Value), 0)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &oldComboboxWindowProc, syscall.NewCallback(comboboxWindowProc))

	retval := &mountedSelectInput{
		Control:       Control{hwnd},
		onChange:      w.OnChange,
		onFocus:       w.OnFocus,
		onBlur:        w.OnBlur,
		longestString: longestString,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedSelectInput struct {
	Control
	onChange func(value int)
	onFocus  func()
	onBlur   func()

	longestString  string
	preferredWidth Length
}

func (w *mountedSelectInput) Layout(bc Constraint) Size {
	return bc.Constrain(w.MinimumSize())
}

func (w *mountedSelectInput) MinimumSize() Size {
	if w.preferredWidth == 0 {
		text, err := syscall.UTF16FromString(w.longestString)
		if err != nil {
			w.preferredWidth = 75 * DIP
		} else {
			width, _ := w.CalcRect(text)
			w.preferredWidth = FromPixelsX(int(width)).Scale(13, 10)
		}
	}

	width := w.preferredWidth
	height := 14 * DIP
	return Size{width, height}
}

func (w *mountedSelectInput) updateProps(data *SelectInput) error {
	// TODO:  Update the items in the combobox
	// TODO:  Update the selection based on Value
	// TODO:  Update the selection based on Unset.

	w.SetDisabled(data.Disabled)
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	// Clear cache
	w.preferredWidth = 0

	return nil
}

func comboboxWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		selectinputGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := selectinputGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := selectinputGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		notification := win.HIWORD(uint32(wParam))
		switch notification {
		case win.CBN_SELCHANGE:
			if w := selectinputGetPtr(hwnd); w.onChange != nil {
				cursel := win.SendMessage(hwnd, win.CB_GETCURSEL, 0, 0)
				w.onChange(int(cursel))
			}
		}
		// defer to old window proc
	}

	return win.CallWindowProc(oldComboboxWindowProc, hwnd, msg, wParam, lParam)
}

func selectinputGetPtr(hwnd win.HWND) *mountedSelectInput {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*mountedSelectInput)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
