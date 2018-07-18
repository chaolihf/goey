package goey

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestImg(t *testing.T) {
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
