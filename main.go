package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"os"
	"runtime/pprof"
	"time"

	"age_of_empires/ecs"
	"age_of_empires/physics"
	"age_of_empires/selection"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/martinlehoux/kagamigo/kcore"
	"golang.org/x/exp/slog"

)

type Action string

const (
	Selecting  Action = "selecting"
	Patrolling Action = "patrolling"
)

var soilColor = color.RGBA{0x60, 0x40, 0x20, 0xff}

type Game struct {
	CurrentAction   Action
	Camera          Camera
	GlobalSelection selection.GlobalSelection
	ResourceAmount  int
	FaceSource      *text.GoTextFaceSource
	UnitBuilder     EntityBuilder
	NowFunc         func() time.Time
	// Components
	Mask []ecs.Mask
	// Entities         []*Entity
	Position         []physics.Point
	RelBounds        []physics.Rectangle
	Image            []*ebiten.Image
	Selectable       []selection.Selectable
	Move             []physics.Move
	Order            []OrderKind
	PatrolOrder      []PatrolOrder
	ResourceGatherer []ResourceGatherer
	ResourceSource   []ResourceSource
	ResourceStorage  []ResourceStorage
	Spawn            []Spawn
}

func (g Game) Now() time.Time {
	if g.NowFunc != nil {
		return g.NowFunc()
	}
	return time.Now()
}

func (g *Game) Append(components Components, mask ecs.Mask) ecs.Entity {
	g.Mask = append(g.Mask, mask)
	g.Position = append(g.Position, components.Position)
	g.RelBounds = append(g.RelBounds, components.RelBounds)
	g.Image = append(g.Image, components.Image)
	g.Selectable = append(g.Selectable, components.Selectable)
	g.Move = append(g.Move, components.Move)
	g.Order = append(g.Order, components.Order)
	g.ResourceGatherer = append(g.ResourceGatherer, components.ResourceGatherer)
	g.ResourceSource = append(g.ResourceSource, components.ResourceSource)
	g.ResourceStorage = append(g.ResourceStorage, components.ResourceStorage)
	g.Spawn = append(g.Spawn, components.Spawn)
	return ecs.Entity(len(g.Mask) - 1)
}

func (g *Game) getMoveMap() physics.MoveMap {
	required := ecs.CM_Position
	blocked := map[physics.Point]bool{}
	for i, mask := range g.Mask {
		if mask&required == required {
			blocked[g.Position[i]] = true
		}
	}
	return physics.MoveMap{Width: WorldWidth, Height: WorldHeight, Blocked: blocked}
}

func (g *Game) entityAt(position physics.Point) ecs.Entity {
	required := ecs.CM_Position
	for i, mask := range g.Mask {
		if mask&required == required && g.Position[i] == position {
			return ecs.Entity(i)
		}
	}
	return -1
}

func getAllStorageDockings(g *Game) []physics.Point {
	storageDockings := []physics.Point{}
	required := ecs.CM_Position | ecs.CM_ResourceStorage
	for i, mask := range g.Mask {
		if mask&required == required {
			position := g.Position[i]
			storageDockings = append(storageDockings, physics.AdjacentPoints(position)...)
		}
	}
	return storageDockings
}

const (
	KeySpawnRequest  = ebiten.KeyS
	KeyPatrolRequest = ebiten.KeyQ
)

func (g *Game) updateSelecting(input Input, moveMap physics.MoveMap) {
	if input.LeftMouseDown {
		g.GlobalSelection.Start(input.Cursor)
		g.GlobalSelection.StartScreen(input.ScreenCursor)
	}
	if input.RightMouseUp {
		destination := input.Cursor.Div(100).Mul(100)
		slog.Info("destination", slog.String("destination", destination.String()))
		for _, i := range g.GlobalSelection.Selected {
			MainAction(g, i, destination, moveMap)
		}
	}
	if input.LeftMouseUp {
		if g.GlobalSelection.IsActive(input.Cursor) {
			if g.GlobalSelection.IsArea(input.Cursor) {
				g.GlobalSelection.SelectMultiple(input.Cursor, g.Mask, g.Selectable, g.Position, g.RelBounds)
			} else {
				g.GlobalSelection.SelectSingle(input.Cursor, g.Mask, g.Selectable, g.Position, g.RelBounds)
			}
		}
	}
	if input.SpawnPressed {
		required := ecs.CM_Spawn
		for _, i := range g.GlobalSelection.Selected {
			if g.Mask[i]&required == required {
				spawn := &g.Spawn[i]
				if g.ResourceAmount >= spawn.UnitResourceCost {
					g.ResourceAmount -= spawn.UnitResourceCost
					spawn.Requests = append(spawn.Requests, SpawnRequest{Start: g.Now()})
				}
			}
		}
	}
}

func (g *Game) UpdateSimulation(input Input) {
	moveMap := g.getMoveMap()

	if input.EscapePressed {
		g.GlobalSelection.Unselect(g.Selectable)
		g.CurrentAction = Selecting
	} else if input.PatrolPressed {
		g.CurrentAction = Patrolling
	}

	switch g.CurrentAction {
	case Selecting:
		g.updateSelecting(input, moveMap)
	case Patrolling:
		if input.RightMouseUp {
			destination := input.Cursor.Div(100).Mul(100)
			StartPatrolSystem(destination, g.Mask, g.Position, g.Move, g.Order, g.PatrolOrder)
			g.CurrentAction = Selecting
		}
	}

	physics.UpdateMoveSystem(moveMap, g.Mask, g.Move, g.Position)
	UpdateSpawnSystem(g, g.Now(), g.Mask, g.Spawn, g.Position)
	UpdateGatherSystem(g, g.Now(), g.getMoveMap(), g.Mask, g.Order, g.Position, g.Move, g.ResourceGatherer)
	UpdatePatrolSystem(g.getMoveMap(), g.Mask, g.Order, g.PatrolOrder, g.Position, g.Move)
}

func (g *Game) Update() error {
	g.Camera.Update(1.0 / 60.0)
	g.UpdateSimulation(CollectInput(&g.Camera))
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
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
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Age of Empire")
	icon, err := os.Open("icon-crop.jpg")
	kcore.Expect(err, "failed to open icon")
	defer icon.Close()
	iconImg, _, err := image.Decode(icon)
	kcore.Expect(err, "failed to decode icon")
	ebiten.SetWindowIcon([]image.Image{iconImg})
	game := &Game{
		ResourceAmount: 150,
		Camera: Camera{
			X:    1000,
			Y:    1000,
			Zoom: 1.0,
		},
	}
	fontData, err := os.ReadFile("fonts/Cinzel-VariableFont_wght.ttf")
	kcore.Expect(err, "failed to read font")
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	kcore.Expect(err, "failed to create font source")
	game.FaceSource = s
	ironMineBuilder := NewEntityBuilder().WithSolid(NewFilledRectangleImage(physics.Point{X: 100, Y: 100}, color.RGBA{0x80, 0x80, 0x80, 0xff})).WithResourceSource(1000).WithSelectable("square", selection.Building)
	game.Append(ironMineBuilder.WithPosition(physics.Point{X: 1000, Y: 1000}).Build())
	townCenterBuilder := NewEntityBuilder().WithSolid(NewFilledRectangleImage(physics.Point{X: 100, Y: 100}, color.RGBA{0x0, 0x0, 0xff, 0xff})).WithResourceStorage().WithSelectable("square", selection.Building).WithSpawn(NewSpawn(50, 1*time.Second))
	townCenter := game.Append(townCenterBuilder.WithPosition(physics.Point{X: 1000, Y: 2000}).Build())
	game.UnitBuilder = NewEntityBuilder().WithSolid(NewFilledCircleImage(100, color.White)).WithSelectable("round", selection.Unit).WithMove().WithOrder().WithResourceGatherer(15)
	for i := 0; i < 0; i++ {
		spawnPosition, _ := physics.Closest(game.Position[townCenter], physics.AdjacentPoints(game.Position[townCenter]), game.getMoveMap())
		game.Append(game.UnitBuilder.WithPosition(spawnPosition).Build())
	}
	game.CurrentAction = Selecting
	ebiten.SetRunnableOnUnfocused(true)
	kcore.Expect(ebiten.RunGame(game), "failed to run game")
}
