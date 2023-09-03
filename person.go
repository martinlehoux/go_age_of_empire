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
	image := ebiten.NewImage(100, 100)
	image.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	selectedImage := ebiten.NewImage(120, 120)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 10)
	selectedImage.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
	selectedImage.DrawImage(image, op)
	return Person{Position: position, size: Point{100, 100}, image: image, selectedImage: selectedImage}
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

func (p *Person) MoveTo(destination Point, blocked map[Point]bool) {
	p.move = NewMove(p.Position, destination, blocked)
}

func (p *Person) Update() {
	p.Position = p.move.Update(p.Position, 10)
}
