// +build go1.12

package goey

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"io"
	"strings"
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
	handle.Set("className", "goey")
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
		Image:  w.propsImage(),
		Width:  w.width,
		Height: w.height,
	}
}

func (w *imgElement) propsImage() image.Image {
	data := w.handle.Get("src").String()

	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[22:]))

	img, err := png.Decode(r)
	if err != nil {
		println("error decoding png: ", err.Error())
	}
	return img
}

func (w *imgElement) updateProps(data *Img) error {
	w.handle.Set("src", imageToAttr(data.Image))
	w.width, w.height = data.Width, data.Height

	return nil
}

func imageToAttr(i image.Image) string {
	ws := bytes.NewBuffer(nil)
	io.WriteString(ws, "data:image/png;base64,")

	// Writing to a memory buffer.  There shouldn't be any errors during the
	// encoding.
	wb := base64.NewEncoder(base64.StdEncoding, ws)
	_ = png.Encode(wb, i)
	wb.Close()

	return ws.String()
}
