// Copyright (c) 2020, Robert W. Johnstone.  All rights reserved.
// Use of this source code is governed by a license, which is specified in the
// LICENSE file.

// Run this package to build the images with the logo for this project.

package main

import (
	"image"
	"image/color"
	"math"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

var (
	// Define our colours
	black  = color.RGBA{0, 0, 0, 0xff}
	accent = color.RGBA{0xcc, 0, 0, 0xff}
)

func main() {
	avatar32()
	avatar64()
	logo256x128()
}

func logo256x128() {
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 256, 128))
	gc := draw2dimg.NewGraphicContext(dest)

	// Initialize our fonts
	draw2d.SetFontFolder("./fonts")
	println(draw2d.FontFileName(draw2d.FontData{
		Name:   "Orbitron",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleBold,
	}))
	println(draw2d.FontFileName(draw2d.FontData{
		Name:   "OpenSans",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleNormal,
	}))

	// Draw the 'G'
	x := 40.0
	y := 40.0
	r := 32.0
	sw := 14.0
	drawG(gc, x, y, r, sw)

	// Rest of the name
	gc.SetFontData(draw2d.FontData{
		Name:   "Orbitron",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleBold,
	})
	gc.SetFontSize(68) // match stroke width to G
	gc.SetFillColor(black)
	right := gc.FillStringAt("oey", x+r+sw/2, y+r+sw/2)
	right += x + r + sw/2
	right -= 6 // Correction for spacing to right of last letter

	// Tag line
	gc.SetFontData(draw2d.FontData{
		Name:   "OpenSans",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleNormal,
	})
	gc.SetFontSize(68)
	left := x - r - sw/2
	tagWidth := stringWidth(gc, "Declarative, cross-platform")
	println((right - left) / tagWidth)
	gc.SetFontSize(68 * (right - left) / tagWidth)
	gc.SetFillColor(accent)
	gc.FillStringAt("Declarative, cross-platform GUIs", left, 128-sw/2)

	// Save to file
	draw2dimg.SaveToPngFile("logo256x128.png", dest)
}

func avatar32() {
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 32, 32))
	gc := draw2dimg.NewGraphicContext(dest)

	drawG(gc, 16, 16, 13, 6)
	draw2dimg.SaveToPngFile("avatar32.png", dest)
}

func avatar64() {
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 64, 64))
	gc := draw2dimg.NewGraphicContext(dest)

	drawG(gc, 32, 32, 26, 12)
	draw2dimg.SaveToPngFile("avatar64.png", dest)
}

func drawG(gc draw2d.GraphicContext, x, y, r, sw float64) {
	println("G:", sw/r)

	gc.SetStrokeColor(black)
	gc.SetLineWidth(sw)
	gc.BeginPath()
	gc.MoveTo(x+r-sw/2, y-r)
	gc.LineTo(x, y-r)
	gc.ArcTo(x, y, r, r, -math.Pi/2, -math.Pi)
	gc.LineTo(x+0.5, y+r)
	gc.Stroke()
	gc.SetFillColor(black)
	gc.BeginPath()
	gc.MoveTo(x, y)
	gc.LineTo(x, y+r)
	gc.ArcTo(x, y, r, r, math.Pi/2, -math.Pi/2)
	gc.Close()
	gc.Fill()
	gc.SetStrokeColor(accent)
	gc.BeginPath()
	gc.MoveTo(x, y+r)
	gc.ArcTo(x+0.5, y+0.5, r-0.5, r-0.5, math.Pi/2, -math.Pi/2)
	gc.LineTo(x+r, y)
	gc.Stroke()
}

func stringWidth(gc draw2d.GraphicContext, s string) float64 {
	left, _, right, _ := gc.GetStringBounds(s)
	return right - left
}
