package main

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"
	"bytes"
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
	UnitBuilder    EntityBuilder
}

// TODO: This could be a maintained index
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

// TODO: This could be a maintained index
func (g *Game) entityAt(position physics.Point) *Entity {
	for _, e := range g.Entities {
		if e.Position.IsEnabled && e.Position.Value == position {
			return e
		}
	}
	return nil
}

func (g *Game) updateSelecting(cursor physics.Point, moveMap physics.MoveMap) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Selection.Start = cursor
		g.Selection.IsActive = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		destination := cursor.Div(100).Mul(100)
		entityAtDestination := g.entityAt(destination)
		slog.Info("destination", slog.String("destination", destination.String()))
		for _, e := range g.Entities {
			e.MainAction(g, destination, entityAtDestination, moveMap)
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
	ironMine := EntityBuilder{}.WithPosition(physics.Point{X: 1000, Y: 1000}).WithImage(NewFilledRectangleImage(physics.Point{X: 100, Y: 100}, color.RGBA{0x80, 0x80, 0x80, 0xff})).WithResourceSource(1000).WithSelection("square").Build()
	game.Entities = append(game.Entities, &ironMine)
	townCenter := EntityBuilder{}.WithPosition(physics.Point{X: 1000, Y: 2000}).WithImage(NewFilledRectangleImage(physics.Point{X: 100, Y: 100}, color.RGBA{0x0, 0x0, 0xff, 0xff})).WithResourceStorage().WithSelection("square").WithSpawn(NewSpawn(50, 5*time.Second)).Build()
	game.Entities = append(game.Entities, &townCenter)
	spawnPosition, _ := game.Closest(townCenter.Position.Value, physics.AdjacentPoints(townCenter.Position.Value))
	game.UnitBuilder = EntityBuilder{}.WithImage(NewFilledCircleImage(100, color.White)).WithSelection("round").WithMove().WithOrder().WithResourceGatherer(15)
	for i := 0; i < 2; i++ {
		spawnPosition, _ = game.Closest(townCenter.Position.Value, physics.AdjacentPoints(townCenter.Position.Value))
		person := game.UnitBuilder.Build()
		person.Position = ecs.C(spawnPosition)
		game.Entities = append(game.Entities, &person)
	}
	game.CurrentAction = Selecting
	ebiten.SetRunnableOnUnfocused(true)
	kcore.Expect(ebiten.RunGame(game), "failed to run game")
}
