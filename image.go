package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewColorImage(size Point, color color.Color) *ebiten.Image {
	image := ebiten.NewImage(size.X, size.Y)
	image.Fill(color)
	return image
}

func NewHaloImage(sourceImage *ebiten.Image, color color.Color, width int) *ebiten.Image {
	image := NewColorImage(sourceImage.Bounds().Size().Add(Point{2 * width, 2 * width}), color)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(width), float64(width))
	image.DrawImage(sourceImage, op)
	return image
}
