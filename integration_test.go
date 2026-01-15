package main

import (
	"testing"
	"time"

	"age_of_empires/physics"
)

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
		WithMove().
		WithOrder().
		WithResourceGatherer(10)
	unitID := game.Append(unitBuilder.Build())

	// Set up simulated time
	simTime := time.Now()
	game.NowFunc = func() time.Time { return simTime }

	// Simulate user selection: select the unit
	game.GlobalSelection.Selected = append(game.GlobalSelection.Selected, unitID)

	// Simulate right-click on mine position to trigger Gather via MainAction
	minePosition := game.Position[mineID]
	MainAction(game, unitID, minePosition, game.getMoveMap())

	initialResources := game.ResourceAmount
	dt := 16 * time.Millisecond
	maxTicks := 100_000

	deposited := false

	for i := 0; i < maxTicks; i++ {
		simTime = simTime.Add(dt)
		game.Update()

		if game.ResourceAmount > initialResources {
			deposited = true
			break
		}
	}

	if !deposited {
		t.Errorf("Unit failed to gather resources within %d ticks. ResourceAmount: %d", maxTicks, game.ResourceAmount)
	}
}
