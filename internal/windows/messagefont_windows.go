package windows

import (
	"fmt"
	"unsafe"

	"github.com/chaolihf/win"
)

const (
	MessageFontScale = 1
)

var (
	hMessageFont win.HFONT
)

func init() {
	// Determine the mssage font
	var ncm win.NONCLIENTMETRICS
	ncm.CbSize = uint32(unsafe.Sizeof(ncm))
	if rc := win.SystemParametersInfo(win.SPI_GETNONCLIENTMETRICS, ncm.CbSize, unsafe.Pointer(&ncm), 0); rc {
		ncm.LfMessageFont.LfHeight = int32(float64(ncm.LfMessageFont.LfHeight) * MessageFontScale)
		ncm.LfMessageFont.LfWidth = int32(float64(ncm.LfMessageFont.LfWidth) * MessageFontScale)
		hMessageFont = win.CreateFontIndirect(&ncm.LfMessageFont)
		if hMessageFont == 0 {
			fmt.Println("Error: failed CreateFontIndirect")
		}
	} else {
		fmt.Println("Error: failed SystemParametersInfo")
	}
}

func MessageFont() win.HFONT {
	return hMessageFont
}
