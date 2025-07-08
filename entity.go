package main

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

type Entity struct {
	Position         ecs.Component[physics.Point]
	Image            ecs.Component[*ebiten.Image]
	Selection        ecs.Component[Selection]
	Move             ecs.Component[physics.Move]
	Order            ecs.Component[Order]
	ResourceGatherer ecs.Component[ResourceGatherer]
	ResourceSource   ecs.Component[ResourceSource]
	ResourceStorage  ecs.Component[ResourceStorage]
	Spawn            ecs.Component[Spawn]
}

func (e Entity) Bounds() physics.Rectangle {
	return physics.Rectangle{
		Min: e.Position.Value,
		Max: e.Position.Value.Add(e.Image.Value.Bounds().Size()),
	}
}
