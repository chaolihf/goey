// This package provides an example application built using the goey package
// that demonstrates using the Image widget.  Clicking the button will cycle
// through images that are a uniform colour, as well as an image of the Go
// mascot.
//
// The management of scrollbars can be tested by using the environment variable
// GOEY_SCROLL.  Allowed values are 0 through 3, which enable no scrollbars,
// the vertical scrollbar, the horizontal scrollbar, or both scrollbars.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"
	"strings"

	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
	"bitbucket.org/rj/goey/windows"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	mainWindow *windows.Window
	clickCount int

	colors = [3]color.RGBA{
		{0xff, 0, 0, 0xff},
		{0, 0xff, 0, 0xff},
		{0, 0, 0xff, 0xff},
	}
	colorNames = [3]string{
		"Red", "Green", "Blue",
	}

	gopher image.Image
)

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		file, err = os.Open("./example/colour/" + filename)
		if err != nil {
			return nil, err
		}
	}
	defer func() {
		// Close of a file with only read permission.  Will not error.
		_ = file.Close()
	}()

	img, _, err := image.Decode(file)
	return img, err
}

func createImage(label string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))

	point := fixed.Point26_6{
		fixed.Int26_6(2 * 64),
		fixed.Int26_6(50 * 64),
	}

	for _, v := range strings.Split(label, "\n") {
		d := &font.Drawer{
			Dst:  img,
			Src:  image.Black,
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(v)

		point.Y += fixed.Int26_6(15 * 64)
	}

	return img
}

func selectImage(index int) (image.Image, string) {
	if clickCount%4 == 3 {
		return gopher, "Image of the Go gopher."
	}

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(img, img.Bounds(), image.NewUniform(colors[index%4]), image.Point{}, draw.Src)
	return img, colorNames[index%4]
}

func main() {
	flag.IntVar(&clickCount, "count", 0, "Initial click count")
	flag.Parse()

	var err error
	gopher, err = loadImage("gopher.png")
	if err != nil {
		fmt.Printf("error: could not open image: %s\n", err)
		gopher = createImage("image not found\n" + err.Error())
	}

	err = loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	// Add the controls
	mw, err := windows.NewWindow("Colour", render())
	if err != nil {
		return err
	}
	mainWindow = mw

	// Set the icon
	img, _ := selectImage(clickCount)
	mw.SetIcon(img)

	return nil
}

func update() {
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	img, _ := selectImage(clickCount)
	mainWindow.SetIcon(img)
}

func render() base.Widget {
	img, description := selectImage(clickCount)

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			AlignMain:  goey.MainCenter,
			AlignCross: goey.CrossCenter,
			Children: []base.Widget{
				&goey.Button{Text: "Change the colour", OnClick: func() {
					clickCount++
					update()
				}},
				&goey.Img{
					Image:  img,
					Width:  (1 * goey.DIP).Scale(img.Bounds().Dx(), 1),
					Height: (1 * goey.DIP).Scale(img.Bounds().Dy(), 1),
				},
				&goey.P{Text: description, Align: goey.JustifyCenter},
			},
		},
	}
}
