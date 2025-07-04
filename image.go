package main

import (
	"age_of_empires/physics"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func NewFilledRectangleImage(size physics.Point, color color.Color) *ebiten.Image {
	image := ebiten.NewImage(size.X, size.Y)
	image.Fill(color)
	return image
}

func NewStrokeRectangleImage(size physics.Point, strokeWidth float32, color color.Color) *ebiten.Image {
	image := ebiten.NewImage(size.X, size.Y)
	vector.StrokeRect(image, 0, 0, float32(size.X), float32(size.Y), float32(strokeWidth), color, true)
	return image
}

func NewFilledCircleImage(width int, color color.Color) *ebiten.Image {
	image := ebiten.NewImage((width), (width))
	vector.DrawFilledCircle(image, float32(width/2), float32(width/2), float32(width/2), color, true)
	return image
}

func NewStrokeCircleImage(width int, strokeWidth int, color color.Color) *ebiten.Image {
	image := ebiten.NewImage((width), (width))
	vector.StrokeCircle(image, float32(width/2), float32(width/2), float32(width/2), float32(strokeWidth), color, true)
	return image
}
