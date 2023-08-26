package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Person struct {
	size          Point
	Position      Point
	IsSelected    bool
	move          Move
	image         *ebiten.Image
	selectedImage *ebiten.Image
}

func NewPerson(position Point) Person {
	image := ebiten.NewImage(10, 10)
	image.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	selectedImage := ebiten.NewImage(12, 12)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(1, 1)
	selectedImage.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
	selectedImage.DrawImage(image, op)
	return Person{Position: position, size: Point{10, 10}, image: image, selectedImage: selectedImage}
}

func (p Person) Image() *ebiten.Image {
	if p.IsSelected {
		return p.selectedImage
	}
	return p.image
}

func (p Person) CollisionBounds() Rectangle {
	return Rectangle{
		p.Position.Sub(p.size.Div(2)),
		p.Position.Add(p.size.Div(2)),
	}
}

func (p *Person) MoveTo(destination Point) {
	p.move = Move{IsActive: true, Destination: destination}
}

func (p *Person) Update() {
	p.Position = p.move.Update(p.Position, 3)
}
