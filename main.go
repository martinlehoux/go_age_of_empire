package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Vec struct {
	X, Y float64
}

type Rectangle struct {
	leftTop     Vec
	rightBottom Vec
}

func NewRectangle(leftTop Vec, rightBottom Vec) Rectangle {
	return Rectangle{leftTop: leftTop, rightBottom: rightBottom}
}

func (r Rectangle) Contains(point Vec) bool {
	return point.X > r.leftTop.X && point.X < r.rightBottom.X && point.Y > r.leftTop.Y && point.Y < r.rightBottom.Y
}

type Person struct {
	size          Vec
	Position      Vec
	IsSelected    bool
	image         *ebiten.Image
	selectedImage *ebiten.Image
}

var personColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
var haloColor = color.RGBA{0xff, 0x00, 0x00, 0xff}

type Tile struct {
	*ebiten.Image
	Size     Vec
	Position Vec
}

func NewTile(size Vec, position Vec, color color.Color) Tile {
	image := ebiten.NewImage(int(size.X), int(size.Y))
	image.Fill(color)
	return Tile{Size: size, Position: position, Image: image}
}

var soilColor = color.RGBA{0x60, 0x40, 0x20, 0xff}
var soil = NewTile(Vec{X: 280, Y: 200}, Vec{X: 20, Y: 20}, soilColor)

func NewPerson(position Vec) Person {
	image := ebiten.NewImage(10, 10)
	image.Fill(personColor)
	selectedImage := ebiten.NewImage(12, 12)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(1, 1)
	selectedImage.Fill(haloColor)
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

func (p *Person) Bounds() Rectangle {
	return NewRectangle(Vec{X: p.Position.X - p.size.X/2, Y: p.Position.Y - p.size.Y/2}, Vec{X: p.Position.X + p.size.X/2, Y: p.Position.Y + p.size.Y/2})
}

func GetDrawImageOptions(size Vec, position Vec) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(position.X-size.X/2, position.Y-size.Y/2)

	return op
}

type Game struct {
	Persons []*Person
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tileClick := Vec{X: float64(x) - soil.Position.X, Y: float64(y) - soil.Position.Y}
		fmt.Println("Click: ", tileClick)
		for _, p := range g.Persons {
			p.IsSelected = false
			fmt.Println("Person: ", p.Position)
			if p.Bounds().Contains(tileClick) {
				fmt.Println("Person selected")
				p.IsSelected = true
			}
		}
	}
	return nil
}

var shallowWaterColor = color.RGBA{0x80, 0xa0, 0xc0, 0xff}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(shallowWaterColor)
	soil.Fill(soilColor)
	for _, p := range g.Persons {
		soil.DrawImage(p.Image(), GetDrawImageOptions(p.Size(), p.Position))
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(soil.Position.X, soil.Position.Y)
	screen.DrawImage(soil.Image, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Age of Empire")
	game := &Game{}
	mainPerson := NewPerson(Vec{X: 120, Y: 80})
	game.Persons = append(game.Persons, &mainPerson)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
