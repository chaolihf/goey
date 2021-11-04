package goeyjs

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"io"
)

// ImageToAttr converts the image to a data URL for inlining in an HTML
// attribute.
func ImageToAttr(i image.Image) string {
	ws := bytes.NewBuffer(nil)
	io.WriteString(ws, "data:image/png;base64,")

	// Writing to a memory buffer.  There shouldn't be any errors during the
	// encoding.
	wb := base64.NewEncoder(base64.StdEncoding, ws)
	_ = png.Encode(wb, i)
	wb.Close()

	return ws.String()
}
