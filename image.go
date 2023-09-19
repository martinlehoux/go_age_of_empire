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
