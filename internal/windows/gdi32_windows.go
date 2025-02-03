package windows

import (
	"errors"
	"image"
	"image/draw"
	"unsafe"

	"github.com/chaolihf/win"
)

func CreateIconFromImage(prop image.Image) (win.HICON, error) {
	// Create a mask for the icon.
	// Currently, we are using a straight white mask, but perhaps this
	// should be a copy of the alpha channel if the source image is
	// RGBA.
	bounds := prop.Bounds()
	imgMask := image.NewGray(prop.Bounds())
	draw.Draw(imgMask, bounds, image.White, image.Point{}, draw.Src)
	hmask, err := CreateBitmapFromImage(imgMask)
	if err != nil {
		return 0, err
	}
	defer win.DeleteObject(win.HGDIOBJ(hmask))

	// Convert the image to a bitmap.
	hbitmap, err := CreateBitmapFromImage(prop)
	if err != nil {
		return 0, err
	}
	defer win.DeleteObject(win.HGDIOBJ(hbitmap))

	// Create the icon
	iconinfo := win.ICONINFO{
		FIcon:    win.TRUE,
		HbmMask:  hmask,
		HbmColor: hbitmap,
	}
	hicon := win.CreateIconIndirect(&iconinfo)
	if hicon == 0 {
		panic("Error in CreateIconIndirect")
	}
	return hicon, nil
}

func checkMemoryBlock(pix []byte) bool {
	start := uintptr(unsafe.Pointer(&pix[0]))
	end := start + uintptr(len(pix))

	// No documentation found, but 0xc000400000 appears to be poison.  Image
	// blocks spanning this address lead to CreateBitmap failures.
	const PoisonAddress = 0xc000400000
	return start <= PoisonAddress && end >= PoisonAddress
}

func imageToBitmapRGBA(img *image.RGBA, pix []byte) (win.HBITMAP, error) {
	// Need to convert RGBA to BGRA.
	for i := 0; i < len(pix); i += 4 {
		// swap the red and green bytes.
		pix[i+0], pix[i+2] = pix[i+2], pix[i+0]
	}

	// Move memory to avoid 'invalid argument' failures.
	if checkMemoryBlock(pix) {
		pix = append([]byte(nil), pix...)
	}

	// The following call also works with 4 channels of 8-bits on a Windows
	// machine, but fails on Wine.  Would like it to work on both to ease
	// CI.
	hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 1, 32, unsafe.Pointer(&pix[0]))
	if hbitmap == 0 {
		return 0, errors.New("call to CreateBitmap failed")
	}
	return hbitmap, nil
}

func CreateBitmapFromImage(prop image.Image) (win.HBITMAP, error) {
	if img, ok := prop.(*image.RGBA); ok {
		// Create a copy of the backing for the pixel data
		buffer := append([]uint8(nil), img.Pix...)
		// Create the bitmap.
		hbitmap, err := imageToBitmapRGBA(img, buffer)
		return hbitmap, err
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)

	// Create the bitmap
	hbitmap, err := imageToBitmapRGBA(img, img.Pix)
	return hbitmap, err
}

func BitmapToImage(hdc win.HDC, hbitmap win.HBITMAP) image.Image {
	bmi := win.BITMAPINFO{}
	bmi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bmi))
	win.GetDIBits(hdc, hbitmap, 0, 0, nil, &bmi, 0)
	if bmi.BmiHeader.BiPlanes == 1 && bmi.BmiHeader.BiBitCount == 32 && bmi.BmiHeader.BiCompression == win.BI_BITFIELDS {
		// Get the pixel data
		buffer := make([]byte, bmi.BmiHeader.BiSizeImage)
		if checkMemoryBlock(buffer) {
			buffer = make([]byte, bmi.BmiHeader.BiSizeImage)
		}
		if cnt := win.GetDIBits(hdc, hbitmap, 0, uint32(bmi.BmiHeader.BiHeight), &buffer[0], &bmi, 0); cnt == 0 {
			return nil
		}

		// Need to convert BGR to RGB
		for i := 0; i < len(buffer); i += 4 {
			buffer[i+0], buffer[i+2] = buffer[i+2], buffer[i+0]
		}
		// In GDI, all bitmaps are bottom up.  We need to reorder the rows
		// before the data can be used for a PNG.
		// TODO:  Combine this pass with the previous.
		stride := int(bmi.BmiHeader.BiWidth) * 4
		for y := 0; y < int(bmi.BmiHeader.BiHeight/2); y++ {
			y2 := int(bmi.BmiHeader.BiHeight) - y - 1
			for x := 0; x < stride; x++ {
				// The stride is always the same as the width?
				buffer[y*stride+x], buffer[y2*stride+x] = buffer[y2*stride+x], buffer[y*stride+x]
			}
		}
		return &image.RGBA{
			Pix:    buffer,
			Stride: int(bmi.BmiHeader.BiWidth * 4),
			Rect:   image.Rect(0, 0, int(bmi.BmiHeader.BiWidth), int(bmi.BmiHeader.BiHeight)),
		}
	}

	return nil
}
