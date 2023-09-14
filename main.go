package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Action string

const (
	Selecting Action = "selecting"
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
func (t Tile) CollisionBounds() Rectangle {
	return Rectangle{t.Position, t.Position.Add(t.Size)}
}

var soilColor = color.RGBA{0x60, 0x40, 0x20, 0xff}
var soil = NewTile(Point{2800, 2000}, Point{200, 200}, soilColor)

type Selection struct {
	IsActive bool
	Start    Point
}

type Game struct {
	Persons       []*Person
	CurrentAction Action
	Selection     Selection
}

func (g *Game) getBlocked() map[Point]bool {
	blocked := map[Point]bool{}
	return blocked
}

func (g *Game) updateSelecting(cursor Point) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Selection.Start = cursor
		g.Selection.IsActive = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		destination := cursor.Sub(soil.Position).Div(100).Mul(100).Add(Point{50, 50})
		for _, p := range g.Persons {
			if p.IsSelected {
				futureBounds := Rectangle{
					destination.Sub(p.size.Div(2)),
					destination.Add(p.size.Div(2)),
				}
				if futureBounds.In(soil.Bounds()) {
					p.MoveTo(destination, g.getBlocked())
				}
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.Selection.IsActive && Distance(g.Selection.Start, cursor) > 10 {
			for _, p := range g.Persons {
				p.IsSelected = false
				selectionBounds := Rectangle{g.Selection.Start.Sub(soil.Position), cursor.Sub(soil.Position)}.Canon()
				if selectionBounds.Overlaps(p.CollisionBounds()) {
					p.IsSelected = true
				}
			}
		} else {
			canBeSelected := true
			for _, p := range g.Persons {
				p.IsSelected = false
				if cursor.Sub(soil.Position).In(p.CollisionBounds()) {
					p.IsSelected = canBeSelected
					canBeSelected = false
				}
			}
		}
		g.Selection.IsActive = false
	}
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	cursor := Point{x, y}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		g.CurrentAction = Selecting
	}
	switch g.CurrentAction {
	case Selecting:
		g.updateSelecting(cursor)
	}
	for _, p := range g.Persons {
		p.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	screenClick := Point{x, y}
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	soil.Fill(soilColor)
	for _, p := range g.Persons {
		bounds := p.Image().Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.Position.X-bounds.Dx()/2), float64(p.Position.Y-bounds.Dy()/2))
		soil.DrawImage(p.Image(), op)
		if p.IsSelected && p.move.IsActive {
			last := p.Position
			for _, point := range p.move.Path {
				vector.StrokeLine(soil.Image, float32(last.X), float32(last.Y), float32(point.X), float32(point.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
				last = point
			}
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(soil.Position.X), float64(soil.Position.Y))
	screen.DrawImage(soil.Image, op)
	if g.Selection.IsActive {
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
	nonPlayerPerson := NewPerson(Point{450, 450})
	game.Persons = append(game.Persons, &mainPerson)
	game.Persons = append(game.Persons, &nonPlayerPerson)
	game.CurrentAction = Selecting
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
