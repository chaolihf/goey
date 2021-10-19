package goey

import (
	"image"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	win2 "bitbucket.org/rj/goey/internal/windows"
	"github.com/lxn/win"
)

func (w *Img) mount(parent base.Control) (base.Element, error) {
	// Create the bitmap
	hbitmap, err := win2.ImageToBitmap(w.Image)
	if err != nil {
		return nil, err
	}

	// Create the control
	const STYLE = win.WS_CHILD | win.WS_VISIBLE | win.SS_BITMAP | win.SS_LEFT
	hwnd, _, err := createControlWindow(0, &staticClassName[0], "", STYLE, parent.HWnd)
	if err != nil {
		return nil, err
	}
	win.SendMessage(hwnd, win2.STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	retval := &imgElement{
		Control: Control{hwnd},
		hbitmap: hbitmap,
		width:   w.Width,
		height:  w.Height,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type imgElement struct {
	Control
	hbitmap win.HBITMAP
	width   base.Length
	height  base.Length
}

func (w *imgElement) Props() base.Widget {
	// Need to recreate the image from the HBITMAP
	hbitmap := win.HBITMAP(win.SendMessage(w.hWnd, win2.STM_GETIMAGE, 0 /*IMAGE_BITMAP*/, 0))
	if hbitmap == 0 {
		return &Img{
			Width:  w.width,
			Height: w.height,
		}
	}

	hdc := win.GetDC(w.hWnd)
	img := win2.BitmapToImage(hdc, hbitmap)
	win.ReleaseDC(w.hWnd, hdc)

	return &Img{
		Image:  img,
		Width:  w.width,
		Height: w.height,
	}
}

func (w *imgElement) SetBounds(bounds base.Rectangle) {
	w.Control.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *imgElement) updateImage(img image.Image) error {
	// Convert the imagme to a bitmap
	hbitmap, err := win2.ImageToBitmap(img)
	if err != nil {
		return err
	}

	// Delete the old bitmap
	if w.hbitmap != 0 {
		win.DeleteObject(win.HGDIOBJ(w.hbitmap))
	}

	// Update the control with the new bitmap
	w.hbitmap = hbitmap
	win.SendMessage(w.hWnd, win2.STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	return nil
}

func (w *imgElement) updateProps(data *Img) error {
	err := w.updateImage(data.Image)
	if err != nil {
		return err
	}

	w.width, w.height = data.Width, data.Height

	return nil
}
