package gtk

// #include "thunks.h"
import "C"
import (
	"image"
	"image/draw"
	"unsafe"
)

func ImageImageData(handle uintptr) []byte {
	length := C.size_t(0)

	data := C.imageImageData(unsafe.Pointer(handle), &length)
	return C.GoBytes(unsafe.Pointer(data), C.int(length))
}

// ImageToRGBA makes a copy of the pixel data, and ensures that the pixel data
// is in the correct format.
func ImageToRGBA(prop image.Image) *image.RGBA {
	// Use existing image if possible
	if img, ok := prop.(*image.RGBA); ok {
		return &image.RGBA{
			Pix:    append([]uint8(nil), img.Pix...),
			Stride: img.Stride,
			Rect:   img.Rect,
		}
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)
	return img
}
