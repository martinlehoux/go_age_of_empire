package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
