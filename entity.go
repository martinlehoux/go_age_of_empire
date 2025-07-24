package main

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"
	"image/color"

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

type EntityBuilder struct {
	entity Entity
}

func NewEntityBuilder() EntityBuilder {
	return EntityBuilder{
		entity: Entity{},
	}
}

func (b EntityBuilder) WithPosition(position physics.Point) EntityBuilder {
	b.entity.Position = ecs.C(position)
	return b
}

func (b EntityBuilder) WithImage(image *ebiten.Image) EntityBuilder {
	b.entity.Image = ecs.C(image)
	return b
}

var red = color.RGBA{0xff, 0x00, 0x00, 0xff}

func (b EntityBuilder) WithSelection(haloKind string) EntityBuilder {
	b.entity.Selection = ecs.C(Selection{
		IsSelected: false,
		Halo:       nil,
	})
	switch haloKind {
	case "round":
		b.entity.Selection.Value.Halo = NewStrokeCircleImage(110, SELECTION_HALO_WIDTH, red)
	case "square":
		b.entity.Selection.Value.Halo = NewStrokeRectangleImage(physics.Point{X: 110, Y: 110}, SELECTION_HALO_WIDTH, red)
	}
	return b
}

func (b EntityBuilder) WithResourceStorage() EntityBuilder {
	b.entity.ResourceStorage = ecs.C(ResourceStorage{})
	return b
}

func (b EntityBuilder) WithResourceSource(amount int) EntityBuilder {
	b.entity.ResourceSource = ecs.C(ResourceSource{Remaining: amount})
	return b
}

func (b EntityBuilder) WithSpawn(spawn Spawn) EntityBuilder {
	b.entity.Spawn = ecs.C(spawn)
	return b
}

func (b EntityBuilder) WithMove() EntityBuilder {
	b.entity.Move = ecs.C(physics.Move{IsActive: false})
	return b
}

func (b EntityBuilder) WithOrder() EntityBuilder {
	b.entity.Order.IsEnabled = true
	return b
}

func (b EntityBuilder) WithResourceGatherer(maxCapacity int) EntityBuilder {
	b.entity.ResourceGatherer = ecs.C(ResourceGatherer{MaxCapacity: maxCapacity})
	return b
}

func (b EntityBuilder) Build() Entity {
	return b.entity
}
