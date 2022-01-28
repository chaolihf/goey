//go:build go1.12
// +build go1.12

package goey

import (
	"encoding/base64"
	"image"
	"image/png"
	"strings"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/js"
)

type imgElement struct {
	Control

	width, height base.Length
}

func (w *Img) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := goeyjs.CreateElement("img", "goey")
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
	w.handle.Set("src", goeyjs.ImageToAttr(data.Image))
	w.width, w.height = data.Width, data.Height

	return nil
}
