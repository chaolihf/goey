package goey

import (
	win2 "bitbucket.org/rj/goey/syscall"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	edit struct {
		className     *uint16
		oldWindowProc uintptr
		emptyString   uint16
	}
)

func init() {
	var err error
	edit.className, err = syscall.UTF16PtrFromString("EDIT")
	if err != nil {
		panic(err)
	}
}

func (w *TextInput) mount(parent Control) (Element, error) {
	text, err := syscall.UTF16PtrFromString(w.Value)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_AUTOHSCROLL)
	if w.Password {
		style = style | win.ES_PASSWORD
	}
	if w.ReadOnly {
		style = style | win.ES_READONLY
	}
	if w.OnEnterKey != nil {
		style = style | win.ES_MULTILINE
	}
	hwnd := win.CreateWindowEx(win.WS_EX_CLIENTEDGE, edit.className, text,
		style,
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

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &edit.oldWindowProc, syscall.NewCallback(textinputWindowProc))

	// Create placeholder, if required.
	if w.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(w.Placeholder)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		win.SendMessage(hwnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	}

	retval := &mountedTextInput{mountedTextInputBase{
		Control:    Control{hwnd},
		onChange:   w.OnChange,
		onFocus:    w.OnFocus,
		onBlur:     w.OnBlur,
		onEnterKey: w.OnEnterKey,
	}}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedTextInputBase struct {
	Control
	onChange   func(value string)
	onFocus    func()
	onBlur     func()
	onEnterKey func(value string)
}

type mountedTextInput struct {
	mountedTextInputBase
}

func (w *mountedTextInput) Props() Widget {
	var buffer [80]uint16
	win.SendMessage(w.hWnd, win.EM_GETCUEBANNER, uintptr(unsafe.Pointer(&buffer[0])), 80)
	ndx := 0
	for i, v := range buffer {
		if v == 0 {
			ndx = i
			break
		}
	}
	placeholder := syscall.UTF16ToString(buffer[:ndx])

	return &TextInput{
		Value:       w.Control.Text(),
		Placeholder: placeholder,
		Disabled:    !win.IsWindowEnabled(w.hWnd),
		Password:    win.SendMessage(w.hWnd, win.EM_GETPASSWORDCHAR, 0, 0) != 0,
		ReadOnly:    (win.GetWindowLong(w.hWnd, win.GWL_STYLE) & win.ES_READONLY) != 0,
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
		OnEnterKey:  w.onEnterKey,
	}
}

func (w *mountedTextInputBase) Layout(bc Constraint) Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(Size{width, height})
}

func (w *mountedTextInputBase) MinIntrinsicHeight(Length) Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP
}

func (w *mountedTextInputBase) MinIntrinsicWidth(Length) Length {
	// TODO
	return 75 * DIP
}

func (w *mountedTextInputBase) updatePlaceholder(text string) error {
	// Update the control
	if text != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(text)
		if err != nil {
			return err
		}

		win.SendMessage(w.hWnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	} else {
		win.SendMessage(w.hWnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(&edit.emptyString)))
	}

	return nil
}

func (w *mountedTextInputBase) updateProps(data *TextInput) error {
	if data.Value != w.Text() {
		w.SetText(data.Value)
	}
	err := w.updatePlaceholder(data.Placeholder)
	if err != nil {
		return err
	}
	w.SetDisabled(data.Disabled)
	if data.Password {
		// TODO:  ???
	} else {
		win.SendMessage(w.hWnd, win.EM_SETPASSWORDCHAR, 0, 0)
	}
	win.SendMessage(w.hWnd, win.EM_SETREADONLY, uintptr(win.BoolToBOOL(data.ReadOnly)), 0)

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	w.onEnterKey = data.OnEnterKey

	return nil
}

func textinputWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		textinputGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := textinputGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := textinputGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_KEYDOWN:
		if wParam == win.VK_RETURN {
			if w := textinputGetPtr(hwnd); w.onEnterKey != nil {
				w.onEnterKey(win2.GetWindowText(hwnd))
				return 0
			}
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		notification := win.HIWORD(uint32(wParam))
		switch notification {
		case win.EN_UPDATE:
			if w := textinputGetPtr(hwnd); w.onChange != nil {
				w.onChange(win2.GetWindowText(hwnd))
			}
		}
		return 0

	}

	return win.CallWindowProc(edit.oldWindowProc, hwnd, msg, wParam, lParam)
}

func textinputGetPtr(hwnd win.HWND) *mountedTextInputBase {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*mountedTextInputBase)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
