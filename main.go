package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

type Selection struct {
	IsActive bool
	Start    Point
}

type Game struct {
	Persons   []*Person
	Selection Selection
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Selection.Start = Point{x, y}
		g.Selection.IsActive = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
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
		tileClick := Point{x - soil.Position.X, y - soil.Position.Y}
		if g.Selection.IsActive && Distance(g.Selection.Start, Point{x, y}) > 10 {
			for _, p := range g.Persons {
				p.IsSelected = false
				selectionBounds := Rectangle{g.Selection.Start.Sub(soil.Position), tileClick}.Canon()
				if selectionBounds.Overlaps(p.CollisionBounds()) {
					p.IsSelected = true
				}
			}
		} else {
			canBeSelected := true
			for _, p := range g.Persons {
				p.IsSelected = false
				if tileClick.In(p.CollisionBounds()) {
					p.IsSelected = canBeSelected
					canBeSelected = false
				}
			}
		}
		g.Selection.IsActive = false

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
		if p.IsSelected && p.move.IsActive {
			vector.StrokeLine(soil.Image, float32(p.Position.X), float32(p.Position.Y), float32(p.move.Destination.X), float32(p.move.Destination.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(soil.Position.X), float64(soil.Position.Y))
	screen.DrawImage(soil.Image, op)
	if g.Selection.IsActive {
		x, y := ebiten.CursorPosition()
		screenClick := Point{x, y}
		vector.StrokeRect(screen, float32(g.Selection.Start.X), float32(g.Selection.Start.Y), float32(screenClick.X-g.Selection.Start.X), float32(screenClick.Y-g.Selection.Start.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 3200, 2400
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Age of Empire")
	game := &Game{}
	mainPerson := NewPerson(Point{soil.Image.Bounds().Dx() / 2, soil.Image.Bounds().Dy() / 2})
	nonPlayerPerson := NewPerson(Point{400, 400})
	game.Persons = append(game.Persons, &mainPerson)
	game.Persons = append(game.Persons, &nonPlayerPerson)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
