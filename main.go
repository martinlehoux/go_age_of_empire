package main

import (
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
var soil = NewTile(Point{2800, 2000}, Point{200, 200}, soilColor)

type Game struct {
	Persons []*Person
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		tileClick := Point{x - soil.Position.X, y - soil.Position.Y}
		for _, p := range g.Persons {
			futureBounds := Rectangle{
				tileClick.Sub(p.size.Div(2)),
				tileClick.Add(p.size.Div(2)),
			}
			if p.IsSelected && futureBounds.In(soil.Bounds()) {
				p.MoveTo(tileClick)
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tileClick := Point{x - soil.Position.X, y - soil.Position.Y}
		for _, p := range g.Persons {
			p.IsSelected = false
			if tileClick.In(p.CollisionBounds()) {
				p.IsSelected = true
			}
		}
	}
	for _, p := range g.Persons {
		p.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	soil.Fill(soilColor)
	for _, p := range g.Persons {
		bounds := p.Image().Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.Position.X-bounds.Dx()/2), float64(p.Position.Y-bounds.Dy()/2))
		soil.DrawImage(p.Image(), op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(soil.Position.X), float64(soil.Position.Y))
	screen.DrawImage(soil.Image, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 3200, 2400
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Age of Empire")
	game := &Game{}
	mainPerson := NewPerson(Point{soil.Image.Bounds().Dx() / 2, soil.Image.Bounds().Dy() / 2})
	game.Persons = append(game.Persons, &mainPerson)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
