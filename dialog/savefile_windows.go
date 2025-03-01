package dialog

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/chaolihf/win"
)

func (m *SaveFile) show() (string, error) {
	title, err := syscall.UTF16PtrFromString(m.title)
	if err != nil {
		return "", err
	}

	ofn := win.OPENFILENAME{
		LStructSize: uint32(unsafe.Sizeof(win.OPENFILENAME{})),
		HwndOwner:   m.owner.HWnd,
		LpstrFilter: buildFilterString(m.filters),
		LpstrTitle:  title,
	}

	filename := [1024]uint16{}
	buffer, err := buildFileString(&ofn, filename[:], m.filename)
	if err != nil {
		return "", err
	}

	rc := win.GetSaveFileName(&ofn)
	if !rc {
		if err := win.CommDlgExtendedError(); err != 0 {
			return "", fmt.Errorf("call to GetOpenFileName failed with code %x", err)
		}
		return "", nil
	}
	return syscall.UTF16ToString(buffer), nil
}
