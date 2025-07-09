package main

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/martinlehoux/kagamigo/kcore"
	"golang.org/x/exp/slog"
	"golang.org/x/image/font/gofont/goregular"
)

type Action string

const (
	Selecting  Action = "selecting"
	Patrolling Action = "patrolling"
)

var soilColor = color.RGBA{0x60, 0x40, 0x20, 0xff}

type GlobalSelection struct {
	IsActive bool
	Start    physics.Point
}

type Game struct {
	Entities       []*Entity
	CurrentAction  Action
	Selection      GlobalSelection
	ResourceAmount int
	FaceSource     *text.GoTextFaceSource
	personImage    *ebiten.Image
	personSelectionHalo *ebiten.Image
}

func DrawMove(screen *ebiten.Image, e *Entity) {
	if e.Move.IsEnabled && e.Position.IsEnabled && e.Selection.IsEnabled {
		if e.Selection.Value.IsSelected && e.Move.Value.IsActive {
			last := e.Position.Value
			dx := +e.Image.Value.Bounds().Dx() / 2
			dy := +e.Image.Value.Bounds().Dy() / 2
			for _, point := range e.Move.Value.Path {
				vector.StrokeLine(screen, float32(last.X+dx), float32(last.Y+dy), float32(point.X+dx), float32(point.Y+dy), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
				last = point
			}
		}
	}
}

func (g *Game) getMoveMap() physics.MoveMap {
	blocked := map[physics.Point]bool{}
	for _, e := range g.Entities {
		if e.Position.IsEnabled {
			blocked[e.Position.Value] = true
		}
	}
	return physics.MoveMap{Width: 3200, Height: 2400, Blocked: blocked}
}

func (g *Game) Closest(from physics.Point, candidates []physics.Point) (physics.Point, int) {
	var closest physics.Point
	closestDistance := math.MaxInt
	moveMap := g.getMoveMap()
	for _, dest := range candidates {
		path, ok := physics.SearchPath(from, dest, moveMap)
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

func (g *Game) updateSelecting(cursor physics.Point, moveMap physics.MoveMap) {
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
					physics.StartMove(&e.Move, e.Position, destination, moveMap)
				}
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.Selection.IsActive && physics.Distance(g.Selection.Start, cursor) > 10 {
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
	if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		for _, e := range g.Entities {
			if !e.Selection.IsEnabled || !e.Selection.Value.IsSelected || !e.Spawn.IsEnabled {
				continue
			}
			e.Spawn.Value.AddRequest(g)
		}
	}
}

func (g *Game) updatePatrolling(cursor physics.Point) {
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
	cursor := physics.Point{X: x, Y: y}
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
		physics.UpdateMove(&e.Move, &e.Position, moveMap)
		e.UpdateOrder(g)
		UpdateSpawn(g, &e.Spawn, e.Position)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	cursor := physics.Point{X: x, Y: y}
	screen.Fill(soilColor)
	for _, e := range g.Entities {
		e.Draw(g, screen)
	}
	for _, e := range g.Entities {
		DrawMove(screen, e)
	}
	if g.Selection.IsActive {
		vector.StrokeRect(screen, float32(g.Selection.Start.X), float32(g.Selection.Start.Y), float32(cursor.X-g.Selection.Start.X), float32(cursor.Y-g.Selection.Start.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
	}
	bannerHeight := float32(200)
	vector.DrawFilledRect(screen, 0, 0, float32(3200), bannerHeight, color.White, true)

	resourceText := fmt.Sprintf("Resources: %d", g.ResourceAmount)
	op := &text.DrawOptions{}
	op.GeoM.Translate(25, 25)
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, resourceText, &text.GoTextFace{
		Source: g.FaceSource,
		Size:   100,
	}, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 3200, 2400
}

func main() {
	for _, arg := range os.Args {
		if arg == "--profile" {
			f, err := os.Create("cpu.prof")
			kcore.Expect(err, "could not create CPU profile")
			defer f.Close()
		    kcore.Expect(pprof.StartCPUProfile(f), "could not start CPU profile")
			defer pprof.StopCPUProfile()
		}
	}
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Age of Empire")
	icon, err := os.Open("icon-crop.jpg")
	kcore.Expect(err, "failed to open icon")
	defer icon.Close()
	iconImg, _, err := image.Decode(icon)
	kcore.Expect(err, "failed to decode icon")
	ebiten.SetWindowIcon([]image.Image{iconImg})
	game := &Game{}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	kcore.Expect(err, "failed to create font source")
	game.FaceSource = s
	ironImage := NewFilledRectangleImage(physics.Point{X: 100, Y: 100}, color.RGBA{0x80, 0x80, 0x80, 0xff})
	ironSelectionHalo := NewStrokeRectangleImage(physics.Point{X: 110, Y: 110}, SELECTION_HALO_WIDTH, color.RGBA{0xff, 0x00, 0x00, 0xff})
	ironMine := Entity{
		Position:       ecs.C(physics.Point{X: 1000, Y: 1000}),
		Image:          ecs.C(ironImage),
		ResourceSource: ecs.C(ResourceSource{Remaining: 1000}),
		Selection:      ecs.C(Selection{IsSelected: false, Halo: ironSelectionHalo}),
	}
	game.Entities = append(game.Entities, &ironMine)
	townCenterImage := NewFilledRectangleImage(physics.Point{X: 100, Y: 100}, color.RGBA{0x00, 0x00, 0xff, 0xff})
	townCenter := Entity{
		Position:        ecs.C(physics.Point{X: 1000, Y: 2000}),
		Image:           ecs.C(townCenterImage),
		ResourceStorage: ecs.C(ResourceStorage{}),
		Selection:       ecs.C(Selection{IsSelected: false, Halo: NewStrokeRectangleImage(physics.Point{X: 110, Y: 110}, SELECTION_HALO_WIDTH, color.RGBA{0xff, 0x00, 0x00, 0xff})}),
		Spawn: ecs.C(NewSpawn(50, 5*time.Second)),
	}
	game.Entities = append(game.Entities, &townCenter)
	var order Order
	game.personImage = NewFilledCircleImage(100, color.RGBA{0xff, 0xff, 0xff, 0xff})
	game.personSelectionHalo = NewStrokeCircleImage(110, SELECTION_HALO_WIDTH, color.RGBA{0xff, 0x00, 0x00, 0xff})
	spawnPosition, _ := game.Closest(townCenter.Position.Value, physics.AdjacentPoints(townCenter.Position.Value))
	person1 := Entity{
		Position:         ecs.C(spawnPosition),
		Image:            ecs.C(game.personImage),
		Selection:        ecs.C(Selection{IsSelected: false, Halo: game.personSelectionHalo}),
		Move:             ecs.C(physics.Move{IsActive: false}),
		Order:            ecs.C(order),
		ResourceGatherer: ecs.C(ResourceGatherer{MaxCapacity: 15}),
	}
	game.Entities = append(game.Entities, &person1)
	spawnPosition, _ = game.Closest(townCenter.Position.Value, physics.AdjacentPoints(townCenter.Position.Value))
	person2 := Entity{
		Position:         ecs.C(spawnPosition),
		Image:            ecs.C(game.personImage),
		Selection:        ecs.C(Selection{IsSelected: false, Halo: game.personSelectionHalo}),
		Move:             ecs.C(physics.Move{IsActive: false}),
		Order:            ecs.C(order),
		ResourceGatherer: ecs.C(ResourceGatherer{MaxCapacity: 15}),
	}
	game.Entities = append(game.Entities, &person2)
	game.CurrentAction = Selecting
	kcore.Expect(ebiten.RunGame(game), "failed to run game")
}
