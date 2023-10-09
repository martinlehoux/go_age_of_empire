package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/slog"
)

type Action string

const (
	Selecting  Action = "selecting"
	Patrolling Action = "patrolling"
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

func (g *Game) Closest(from Point, candidates []Point) (Point, int) {
	var closest Point
	closestDistance := math.MaxInt
	moveMap := g.getMoveMap()
	for _, dest := range candidates {
		path, ok := SearchPath(from, dest, moveMap)
		if !ok {
			continue
		}
		if len(path) < closestDistance {
			closest = dest
			closestDistance = len(path)
		}
	}
	return closest, closestDistance
}

func (g *Game) updateSelecting(cursor Point, moveMap MoveMap) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Selection.Start = cursor
		g.Selection.IsActive = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		destination := cursor.Div(100).Mul(100)
		var entityAtDestination *Entity
		for _, e := range g.Entities {
			if e.Position.IsEnabled && e.Position.Value == destination {
				entityAtDestination = e
				break
			}
		}
		slog.Info("destination", slog.String("destination", destination.String()))
		for _, e := range g.Entities {
			if e.Selection.IsEnabled && e.Selection.Value.IsSelected {
				if entityAtDestination != nil && entityAtDestination.ResourceSource.IsEnabled {
					Gather(e, entityAtDestination, g)
				} else {
					e.StartMove(destination, moveMap)
				}
			}
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

func (g *Game) updatePatrolling(cursor Point) {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		destination := cursor.Div(100).Mul(100)
		for _, e := range g.Entities {
			Patrol(e, destination)
		}
		g.CurrentAction = Selecting
	}
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	cursor := Point{x, y}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		g.CurrentAction = Selecting
		slog.Info("selecting action")
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyA) { // Should be Q
		g.CurrentAction = Patrolling
		slog.Info("patrolling action")
	}
	moveMap := g.getMoveMap()
	switch g.CurrentAction {
	case Selecting:
		g.updateSelecting(cursor, moveMap)
	case Patrolling:
		g.updatePatrolling(cursor)
	}
	for _, e := range g.Entities {
		e.UpdateMove(moveMap)
		e.UpdateOrder(g)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	cursor := Point{x, y}
	screen.Fill(soilColor)
	for _, e := range g.Entities {
		Draw(screen, e)
	}
	for _, e := range g.Entities {
		DrawSelection(screen, e)
	}
	for _, e := range g.Entities {
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
	icon, err := os.Open("icon-crop.jpg")
	if err != nil {
		panic(err)
	}
	defer icon.Close()
	iconImg, _, err := image.Decode(icon)
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowIcon([]image.Image{iconImg})
	game := &Game{}
	ironImage := NewFilledRectangleImage(Point{100, 100}, color.RGBA{0x80, 0x80, 0x80, 0xff})
	ironMine := Entity{
		Position:       C(Point{1000, 1000}),
		Image:          C(ironImage),
		ResourceSource: C(ResourceSource{Remaining: 1000}),
	}
	game.Entities = append(game.Entities, &ironMine)
	storageImage := NewFilledRectangleImage(Point{100, 100}, color.RGBA{0x00, 0x00, 0xff, 0xff})
	storage := Entity{
		Position:        C(Point{1000, 2000}),
		Image:           C(storageImage),
		ResourceStorage: C(ResourceStorage{}),
	}
	game.Entities = append(game.Entities, &storage)
	var order Order
	personImage := NewFilledCircleImage(100, color.RGBA{0xff, 0xff, 0xff, 0xff})
	personSelectionHalo := NewStrokeCircleImage(110, SELECTION_HALO_WIDTH, color.RGBA{0xff, 0x00, 0x00, 0xff})
	person1 := Entity{
		Position:         C(Point{2000, 2000}),
		Image:            C(personImage),
		Selection:        C(Selection{IsSelected: false, Halo: personSelectionHalo}),
		Move:             C(Move{IsActive: false}),
		Order:            C(order),
		ResourceGatherer: C(ResourceGatherer{MaxCapacity: 15}),
	}
	game.Entities = append(game.Entities, &person1)
	person2 := Entity{
		Position:         C(Point{2200, 2200}),
		Image:            C(personImage),
		Selection:        C(Selection{IsSelected: false, Halo: personSelectionHalo}),
		Move:             C(Move{IsActive: false}),
		Order:            C(order),
		ResourceGatherer: C(ResourceGatherer{MaxCapacity: 15}),
	}
	game.Entities = append(game.Entities, &person2)
	game.CurrentAction = Selecting
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
