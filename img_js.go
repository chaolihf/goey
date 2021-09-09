package goey

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"io"
	"syscall/js"

	"bitbucket.org/rj/goey/base"
)

type imgElement struct {
	Control

	width, height base.Length
}

func (w *Img) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := js.Global().Get("document").Call("createElement", "img")
	handle.Get("style").Set("position", "absolute")
	parent.Handle.Call("appendChild", handle)

	// Create the element
	retval := &imgElement{
		Control: Control{handle},
	}
	retval.updateProps(w)

	return retval, nil
}

func (w *imgElement) Props() base.Widget {
	return &Img{
		Image:  w.PropsImage(),
		Width:  w.width,
		Height: w.height,
	}
}

func (w *imgElement) PropsImage() image.Image {
	return nil
}

func (w *imgElement) updateProps(data *Img) error {
	w.handle.Set("src", imageToAttr(data.Image))
	w.width, w.height = data.Width, data.Height

	return nil
}

func imageToAttr(i image.Image) string {
	w := bytes.NewBuffer(nil)
	io.WriteString(w, "data:image/png;base64,")

	// Writing to a memory buffer.  There shouldn't be any errors during the
	// encoding.
	_ = png.Encode(base64.NewEncoder(base64.StdEncoding, w), i)

	return w.String()
}
