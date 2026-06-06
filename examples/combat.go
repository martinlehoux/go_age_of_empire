package examples

import (
	"image/color"

	"age_of_empires/engine"
	"age_of_empires/physics"
	"age_of_empires/selection"
)

func init() {
	Registry["combat"] = Entry{
		Name:        "Combat",
		Description: "Two groups of units facing each other — ready for battle",
		Build:       buildCombat,
	}
}

func buildCombat(g *engine.Game, deps BuilderDeps) {
	g.ResourceAmount = 500
	g.Camera = engine.Camera{X: 800, Y: 1000, Zoom: 1.0}
	g.CurrentAction = engine.Selecting
	g.UnitBuilder = deps.UnitBuilder

	// Red team (left)
	redBuilder := engine.NewEntityBuilder().
		WithSolid(engine.NewFilledCircleImage(100, color.RGBA{0xff, 0x00, 0x00, 0xff})).
		WithSelectable("round", selection.Unit).
		WithMove().WithOrder().WithResourceGatherer(15)
	for x := 0; x < 5; x++ {
		g.Append(redBuilder.WithPosition(physics.Point{X: 200 + x*120, Y: 1000}).Build())
	}

	// Blue team (right)
	blueBuilder := engine.NewEntityBuilder().
		WithSolid(engine.NewFilledCircleImage(100, color.RGBA{0x00, 0x00, 0xff, 0xff})).
		WithSelectable("round", selection.Unit).
		WithMove().WithOrder().WithResourceGatherer(15)
	for x := 0; x < 5; x++ {
		g.Append(blueBuilder.WithPosition(physics.Point{X: 1000 + x*120, Y: 1000}).Build())
	}
}
