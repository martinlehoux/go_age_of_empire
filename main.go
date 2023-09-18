package main

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/slog"
)

type Action string

const (
	Selecting Action = "selecting"
)

var soilColor = color.RGBA{0x60, 0x40, 0x20, 0xff}

type GlobalSelection struct {
	IsActive bool
	Start    Point
}

type Game struct {
	Entities      []*Entity
	CurrentAction Action
	Selection     GlobalSelection
}

func (g *Game) getMoveMap() MoveMap {
	blocked := map[Point]bool{}
	for _, e := range g.Entities {
		if e.Position.IsEnabled {
			blocked[e.Position.Value] = true
		}
	}
	return MoveMap{Width: 3200, Height: 2400, Blocked: blocked}
}

func (g *Game) updateSelecting(cursor Point) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Selection.Start = cursor
		g.Selection.IsActive = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		destination := cursor.Div(100).Mul(100)
		slog.Info("destination", slog.String("destination", destination.String()))
		moveMap := g.getMoveMap()
		for _, e := range g.Entities {
			e.StartMove(destination, moveMap)
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.Selection.IsActive && Distance(g.Selection.Start, cursor) > 10 {
			for _, e := range g.Entities {
				e.SelectMultiple(cursor, g.Selection)
			}
		} else {
			canBeSelected := true
			for _, e := range g.Entities {
				if e.SelectSingle(cursor, canBeSelected) {
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
	for _, e := range g.Entities {
		e.UpdateMove()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	cursor := Point{x, y}
	screen.Fill(soilColor)
	for _, e := range g.Entities {
		Draw(screen, e)
		DrawMove(screen, e)
	}
	if g.Selection.IsActive {
		vector.StrokeRect(screen, float32(g.Selection.Start.X), float32(g.Selection.Start.Y), float32(cursor.X-g.Selection.Start.X), float32(cursor.Y-g.Selection.Start.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 3200, 2400
}

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Age of Empire")
	game := &Game{}
	ironImage := NewColorImage(Point{100, 100}, color.RGBA{0x80, 0x80, 0x80, 0xff})
	ironSelectedImage := NewHaloImage(ironImage, color.RGBA{0xff, 0xff, 0xff, 0xff}, 10)
	ironMine := Entity{
		Position:  C(Point{1000, 1000}),
		Image:     C(ironImage),
		Selection: C(Selection{IsSelected: false, SelectedImage: ironSelectedImage}),
	}
	game.Entities = append(game.Entities, &ironMine)
	personImage := NewColorImage(Point{100, 100}, color.RGBA{0xff, 0xff, 0xff, 0xff})
	personSelectedImage := NewHaloImage(personImage, color.RGBA{0xff, 0x00, 0x00, 0xff}, 10)
	alterPerson := Entity{
		Position:  C(Point{2000, 2000}),
		Image:     C(personImage),
		Selection: C(Selection{IsSelected: false, SelectedImage: personSelectedImage}),
		Move:      C(Move{IsActive: false}),
	}
	game.Entities = append(game.Entities, &alterPerson)
	game.CurrentAction = Selecting
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
