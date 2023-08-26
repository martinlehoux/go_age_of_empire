package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Person struct {
	size          Vec
	Position      Vec
	IsSelected    bool
	image         *ebiten.Image
	selectedImage *ebiten.Image
}

func NewPerson(position Vec) Person {
	image := ebiten.NewImage(10, 10)
	image.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	selectedImage := ebiten.NewImage(12, 12)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(1, 1)
	selectedImage.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
	selectedImage.DrawImage(image, op)
	return Person{Position: position, size: Vec{X: 10, Y: 10}, image: image, selectedImage: selectedImage}
}

func (p Person) Image() *ebiten.Image {
	if p.IsSelected {
		return p.selectedImage
	}
	return p.image
}

func (p Person) Size() Vec {
	if p.IsSelected {
		return Vec{X: p.size.X + 2, Y: p.size.Y + 2}
	} else {
		return p.size
	}
}

func (p Person) Bounds() Rectangle {
	return Rectangle{Vec{p.Position.X - p.size.X/2, p.Position.Y - p.size.Y/2}, Vec{p.Position.X + p.size.X/2, p.Position.Y + p.size.Y/2}}
}
