package examples

import (
	"image/color"
	"time"

	"age_of_empires/engine"
	"age_of_empires/physics"
	"age_of_empires/selection"
)

func init() {
	Registry["default"] = Entry{
		Name:        "Default",
		Description: "Free play with a town center, iron mine, and gatherer units",
		Build:       buildDefault,
	}
}

func buildDefault(g *engine.Game, deps BuilderDeps) {
	g.ResourceAmount = 150
	g.Camera = engine.Camera{X: 1000, Y: 1000, Zoom: 1.0}
	g.CurrentAction = engine.Selecting
	g.UnitBuilder = deps.UnitBuilder

	ironMineBuilder := engine.NewEntityBuilder().
		WithSolid(engine.NewFilledRectangleImage(physics.Point{X: 100, Y: 100}, color.RGBA{0x80, 0x80, 0x80, 0xff})).
		WithResourceSource(1000).
		WithSelectable("square", selection.Building)
	g.Append(ironMineBuilder.WithPosition(physics.Point{X: 1000, Y: 1000}).Build())

	townCenterBuilder := engine.NewEntityBuilder().
		WithSolid(engine.NewFilledRectangleImage(physics.Point{X: 100, Y: 100}, color.RGBA{0x0, 0x0, 0xff, 0xff})).
		WithResourceStorage().
		WithSelectable("square", selection.Building).
		WithSpawn(engine.NewSpawn(50, 1*time.Second))
	townCenter := g.Append(townCenterBuilder.WithPosition(physics.Point{X: 1000, Y: 2000}).Build())

	_ = townCenter // suppress unused warning; future spawn loop will use it
}
