package dialog

import (
	"syscall"

	"github.com/chaolihf/win"
)

func (m *Message) show() error {
	text, err := syscall.UTF16PtrFromString(m.text)
	if err != nil {
		return err
	}
	title, err := syscall.UTF16PtrFromString(m.title)
	if err != nil {
		return err
	}

	rc := win.MessageBox(m.owner.HWnd, text, title, uint32(m.icon))
	if rc == 0 {
		return syscall.GetLastError()
	}
	return nil
}

func (m *Message) withError() {
	m.icon |= win.MB_ICONERROR
}

func (m *Message) withWarn() {
	m.icon |= win.MB_ICONWARNING
}

func (m *Message) withInfo() {
	m.icon |= win.MB_ICONINFORMATION
}
