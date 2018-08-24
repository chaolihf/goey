package goey

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestImgCreate(t *testing.T) {
	bounds := image.Rect(0, 0, 92, 92)
	images := []*image.RGBA{image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds)}
	draw.Draw(images[0], bounds, image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	draw.Draw(images[1], bounds, image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(images[2], bounds, image.NewUniform(color.RGBA{0, 255, 255, 255}), image.Point{}, draw.Src)

	testingRenderWidgets(t,
		&Img{Image: images[0], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[1]},
		&Img{Image: images[2]},
	)
}

func TestImgClose(t *testing.T) {
	bounds := image.Rect(0, 0, 92, 92)
	images := []*image.RGBA{image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds)}
	draw.Draw(images[0], bounds, image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	draw.Draw(images[1], bounds, image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(images[2], bounds, image.NewUniform(color.RGBA{0, 255, 255, 255}), image.Point{}, draw.Src)

	testingRenderWidgets(t,
		&Img{Image: images[0], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[1], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[2]},
	)
}

func TestImgUpdate(t *testing.T) {
	bounds := image.Rect(0, 0, 92, 92)
	images := []*image.RGBA{image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds)}
	draw.Draw(images[0], bounds, image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	draw.Draw(images[1], bounds, image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(images[2], bounds, image.NewUniform(color.RGBA{0, 255, 255, 255}), image.Point{}, draw.Src)

	testingUpdateWidgets(t, []base.Widget{
		&Img{Image: images[0], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[1], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[2], Width: 100 * DIP, Height: 10 * DIP},
	}, []base.Widget{
		&Img{Image: images[2], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[1], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[0], Width: 100 * DIP, Height: 10 * DIP},
	})
}

func TestImgUpdateDimensions(t *testing.T) {
	img1 := image.RGBA{Rect: image.Rect(0, 0, 92, 92)}

	cases := []struct {
		width  base.Length
		height base.Length
		img    image.Image
		out    base.Size
	}{
		{10 * DIP, 15 * DIP, &img1, base.Size{10 * DIP, 15 * DIP}},
		{0, 0, &img1, base.Size{1 * Inch, 1 * Inch}},
		{2 * Inch, 0, &img1, base.Size{2 * Inch, 2 * Inch}},
		{0, 2 * Inch, &img1, base.Size{2 * Inch, 2 * Inch}},
	}

	for i, v := range cases {
		widget := Img{
			Width:  v.width,
			Height: v.height,
			Image:  v.img,
		}

		widget.UpdateDimensions()
		if widget.Height != v.out.Height {
			t.Errorf("Case %d:  Failed to update height, got %v, want %v", i, widget.Height, v.out.Height)
		}
		if widget.Width != v.out.Width {
			t.Errorf("Case %d:  Failed to update width, got %v, want %v", i, widget.Width, v.out.Width)
		}
	}
}
