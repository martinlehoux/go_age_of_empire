package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Tile struct {
	*ebiten.Image
	Size     Point
	Position Point
}

func NewTile(size Point, position Point, color color.Color) Tile {
	image := ebiten.NewImage(int(size.X), int(size.Y))
	image.Fill(color)
	return Tile{Size: size, Position: position, Image: image}
}

var soilColor = color.RGBA{0x60, 0x40, 0x20, 0xff}
var soil = NewTile(Point{280, 200}, Point{20, 20}, soilColor)

type Game struct {
	Persons []*Person
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tileClick := Point{x - soil.Position.X, y - soil.Position.Y}
		fmt.Println("Click: ", tileClick)
		for _, p := range g.Persons {
			p.IsSelected = false
			fmt.Println("Person: ", p.Position)
			if tileClick.In(p.CollisionBounds()) {
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
		bounds := p.Image().Bounds()
		fmt.Println("Person bounds: ", bounds)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.Position.X-bounds.Dx()/2), float64(p.Position.Y-bounds.Dy()/2))
		fmt.Println("Person position:", p.Position, "Draw geom: ", op.GeoM)
		soil.DrawImage(p.Image(), op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(soil.Position.X), float64(soil.Position.Y))
	screen.DrawImage(soil.Image, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Age of Empire")
	game := &Game{}
	mainPerson := NewPerson(Point{soil.Image.Bounds().Dx() / 2, soil.Image.Bounds().Dy() / 2})
	game.Persons = append(game.Persons, &mainPerson)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
