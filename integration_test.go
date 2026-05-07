package main

import (
	"image/color"
	"testing"
	"time"

	"age_of_empires/physics"
	"age_of_empires/selection"

	"github.com/stretchr/testify/assert"
)

func TestCameraScreenToWorldRoundTrip(t *testing.T) {
	zooms := []float64{0.25, 0.5, 1.0, 2.0, 3.0}
	points := []physics.Point{
		{X: 1000, Y: 1000},
		{X: 0, Y: 0},
		{X: 3200, Y: 2400},
		{X: 1600, Y: 1200},
	}
	for _, zoom := range zooms {
		cam := Camera{X: 1600, Y: 1200, Zoom: zoom}
		for _, p := range points {
			sx, sy := cam.WorldToScreen(p)
			got := cam.ScreenToWorld(int(sx), int(sy))
			assert.Equal(t, p, got, "zoom=%.2f point=%v", zoom, p)
		}
	}
}

func TestCameraClampsToBounds(t *testing.T) {
	cases := []struct{ x, y float64 }{
		{-99999, -99999},
		{99999, 99999},
		{-99999, 99999},
		{99999, -99999},
	}
	for _, c := range cases {
		cam := Camera{X: c.x, Y: c.y, Zoom: 1.0}
		cam.clampToBounds()
		halfW := (ScreenWidth / 2.0) / cam.Zoom
		halfH := (ScreenHeight / 2.0) / cam.Zoom
		assert.GreaterOrEqual(t, cam.X, halfW, "X=%v should be >= halfW", c.x)
		assert.LessOrEqual(t, cam.X, float64(WorldWidth)-halfW, "X=%v should be <= maxX", c.x)
		assert.GreaterOrEqual(t, cam.Y, halfH, "Y=%v should be >= halfH", c.y)
		assert.LessOrEqual(t, cam.Y, float64(WorldHeight)-halfH, "Y=%v should be <= maxY", c.y)
	}
}

func TestCameraZoomTowardCursor(t *testing.T) {
	zooms := []struct{ before, after float64 }{
		{1.0, 2.0},
		{1.0, 0.5},
		{2.0, 3.0},
	}
	cursorSX, cursorSY := 400, 300
	for _, z := range zooms {
		cam := Camera{X: 1600, Y: 1200, Zoom: z.before}
		worldBefore := cam.ScreenToWorld(cursorSX, cursorSY)

		// replicate the zoom-toward-cursor logic from Camera.Update
		wx0 := float64(cursorSX-ScreenWidth/2)/cam.Zoom + cam.X
		wy0 := float64(cursorSY-ScreenHeight/2)/cam.Zoom + cam.Y
		cam.Zoom = z.after
		wx1 := float64(cursorSX-ScreenWidth/2)/cam.Zoom + cam.X
		wy1 := float64(cursorSY-ScreenHeight/2)/cam.Zoom + cam.Y
		cam.X -= wx1 - wx0
		cam.Y -= wy1 - wy0

		worldAfter := cam.ScreenToWorld(cursorSX, cursorSY)
		assert.Equal(t, worldBefore, worldAfter, "zoom %.1f→%.1f: world point under cursor should be stable", z.before, z.after)
	}
}

func TestGatherLoop(t *testing.T) {
	// Initialize Game
	game := &Game{}

	// Town Center (Storage) at -100,0 (different grid cell from unit at 0,0)
	tcBuilder := NewEntityBuilder().
		WithPosition(physics.Point{X: -100, Y: 0}).
		WithResourceStorage()
	_ = game.Append(tcBuilder.Build())

	// Mine (Source) at 0, 200
	mineBuilder := NewEntityBuilder().
		WithPosition(physics.Point{X: 0, Y: 200}).
		WithResourceSource(100)
	mineID := game.Append(mineBuilder.Build())

	// Unit (Gatherer) at 0, 0
	unitBuilder := NewEntityBuilder().
		WithPosition(physics.Point{X: 0, Y: 0}).
		WithSolid(NewFilledCircleImage(100, color.White)).
		WithSelectable("round", selection.Unit).
		WithMove().
		WithOrder().
		WithResourceGatherer(10)
	game.Append(unitBuilder.Build())

	// Set up simulated time
	simTime := time.Now()
	game.NowFunc = func() time.Time { return simTime }

	// Click unit to select, then right-click mine — same gesture as in-game
	game.UpdateSimulation(Input{EscapePressed: true})
	game.UpdateSimulation(Input{Cursor: physics.Point{X: 50, Y: 50}, LeftMouseDown: true})
	game.UpdateSimulation(Input{Cursor: physics.Point{X: 50, Y: 50}, LeftMouseUp: true})
	game.UpdateSimulation(Input{Cursor: game.Position[mineID], RightMouseUp: true})

	initialResources := game.ResourceAmount
	dt := 16 * time.Millisecond
	maxTicks := 100_000

	deposited := false

	for i := 0; i < maxTicks; i++ {
		simTime = simTime.Add(dt)
		game.UpdateSimulation(Input{})

		if game.ResourceAmount > initialResources {
			deposited = true
			break
		}
	}

	if !deposited {
		t.Errorf("Unit failed to gather resources within %d ticks. ResourceAmount: %d", maxTicks, game.ResourceAmount)
	}
}

func TestTwoEntitiesMovingToSameCellDontOverlap(t *testing.T) {
	// Two units close together both ordered to the same cell should not overlap.
	game := &Game{}

	target := physics.Point{X: 500, Y: 0}
	unitImg := NewFilledCircleImage(100, color.White)

	unitA := game.Append(NewEntityBuilder().
		WithPosition(physics.Point{X: 200, Y: 0}).
		WithSolid(unitImg).WithSelectable("round", selection.Unit).
		WithMove().WithOrder().Build())
	unitB := game.Append(NewEntityBuilder().
		WithPosition(physics.Point{X: 100, Y: 0}).
		WithSolid(unitImg).WithSelectable("round", selection.Unit).
		WithMove().WithOrder().Build())

	// Area-select both units then right-click target — same gesture as in-game
	game.UpdateSimulation(Input{EscapePressed: true})
	game.UpdateSimulation(Input{Cursor: physics.Point{X: 50, Y: -50}, LeftMouseDown: true})
	game.UpdateSimulation(Input{Cursor: physics.Point{X: 350, Y: 50}, LeftMouseUp: true})
	game.UpdateSimulation(Input{Cursor: target, RightMouseUp: true})

	// MOVE_SPEED=10, cell=100 → 10 ticks per cell. B detours around A → ~60 ticks.
	for range 60 {
		game.UpdateSimulation(Input{})
	}

	assert.False(t, game.Move[unitA].IsActive, "Entity A should have stopped moving")
	assert.False(t, game.Move[unitB].IsActive, "Entity B should have stopped moving")
	assert.NotEqual(t, game.Position[unitA], game.Position[unitB],
		"Entity A and B should NOT overlap: A at %s, B at %s", game.Position[unitA], game.Position[unitB])
}

func TestUnitDoesNotStopWhenDestinationOccupied(t *testing.T) {
	// Unit B is far away, unit A is close. Both ordered to the same cell.
	// A arrives first and settles. B should repath to an adjacent free cell
	// rather than stopping mid-journey.
	game := &Game{}
	game.CurrentAction = Selecting
	unitImg := NewFilledCircleImage(100, color.White)

	target := physics.Point{X: 500, Y: 0}

	unitA := game.Append(NewEntityBuilder().
		WithPosition(physics.Point{X: 400, Y: 0}).
		WithSolid(unitImg).WithSelectable("round", selection.Unit).
		WithMove().WithOrder().Build())
	unitB := game.Append(NewEntityBuilder().
		WithPosition(physics.Point{X: 0, Y: 0}).
		WithSolid(unitImg).WithSelectable("round", selection.Unit).
		WithMove().WithOrder().Build())

	game.UpdateSimulation(Input{EscapePressed: true})
	game.UpdateSimulation(Input{Cursor: physics.Point{X: 50, Y: -50}, LeftMouseDown: true})
	game.UpdateSimulation(Input{Cursor: physics.Point{X: 450, Y: 50}, LeftMouseUp: true})
	game.UpdateSimulation(Input{Cursor: target, RightMouseUp: true})

	// Run long enough for both units to settle (B travels ~5 cells)
	for range 100 {
		game.UpdateSimulation(Input{})
	}

	assert.False(t, game.Move[unitA].IsActive, "unit A should have stopped")
	assert.False(t, game.Move[unitB].IsActive, "unit B should have stopped near target, not mid-journey")
	assert.NotEqual(t, game.Position[unitA], game.Position[unitB], "units should not overlap")
	// B should be adjacent to target, not stuck far away
	assert.LessOrEqual(t, physics.Distance(game.Position[unitB], target), 150.0,
		"unit B should be adjacent to target, got %s", game.Position[unitB])
}
