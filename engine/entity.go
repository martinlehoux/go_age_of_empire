package engine

import (
	"image/color"

	"age_of_empires/ecs"
	"age_of_empires/physics"
	"age_of_empires/selection"

	"github.com/hajimehoshi/ebiten/v2"
)

type Components struct {
	Position         physics.Point
	RelBounds        physics.Rectangle
	Image            *ebiten.Image
	Selectable        selection.Selectable
	Move             physics.Move
	Order            OrderKind
	ResourceGatherer ResourceGatherer
	ResourceSource   ResourceSource
	ResourceStorage  ResourceStorage
	Spawn            Spawn
}

type EntityBuilder struct {
	mask       ecs.Mask
	components Components
}

func NewEntityBuilder() EntityBuilder {
	return EntityBuilder{
		components: Components{},
	}
}

func (b EntityBuilder) WithPosition(position physics.Point) EntityBuilder {
	b.components.Position = position
	b.mask |= ecs.CM_Position
	return b
}

func (b EntityBuilder) WithSolid(image *ebiten.Image) EntityBuilder {
	b.components.Image = image
	b.mask |= ecs.CM_Image
	b.components.RelBounds = image.Bounds()
	b.mask |= ecs.CM_RelBounds
	return b
}

var Red = color.RGBA{0xff, 0x00, 0x00, 0xff}

func (b EntityBuilder) WithSelectable(haloKind string, priority selection.Priority) EntityBuilder {
	b.components.Selectable = selection.Selectable{
		IsSelected: false,
		Halo:       nil,
		Priority:   priority,
	}
	b.mask |= ecs.CM_Selectable
	switch haloKind {
	case "round":
		b.components.Selectable.Halo = NewStrokeCircleImage(110, selection.HaloWidth, Red)
	case "square":
		b.components.Selectable.Halo = NewStrokeRectangleImage(physics.Point{X: 110, Y: 110}, selection.HaloWidth, Red)
	}
	return b
}

func (b EntityBuilder) WithResourceStorage() EntityBuilder {
	b.components.ResourceStorage = ResourceStorage{}
	b.mask |= ecs.CM_ResourceStorage
	return b
}

func (b EntityBuilder) WithResourceSource(amount int) EntityBuilder {
	b.components.ResourceSource = ResourceSource{Remaining: amount}
	b.mask |= ecs.CM_ResourceSource
	return b
}

func (b EntityBuilder) WithSpawn(spawn Spawn) EntityBuilder {
	b.components.Spawn = spawn
	b.mask |= ecs.CM_Spawn
	return b
}

func (b EntityBuilder) WithMove() EntityBuilder {
	b.components.Move = physics.Move{IsActive: false}
	b.mask |= ecs.CM_Move
	return b
}

func (b EntityBuilder) WithOrder() EntityBuilder {
	b.components.Order = OrderKindNone
	b.mask |= ecs.CM_Order
	return b
}

func (b EntityBuilder) WithResourceGatherer(maxCapacity int) EntityBuilder {
	b.components.ResourceGatherer = ResourceGatherer{MaxCapacity: maxCapacity}
	b.mask |= ecs.CM_ResourceGatherer
	return b
}

func (b EntityBuilder) Build() (Components, ecs.Mask) {
	return b.components, b.mask
}
