package main

import "github.com/hajimehoshi/ebiten/v2"

type Component[T any] struct {
	IsEnabled bool
	Value     T
}

func C[T any](t T) Component[T] {
	return Component[T]{IsEnabled: true, Value: t}
}

type Entity struct {
	Position  Component[Point]
	Image     Component[*ebiten.Image]
	Selection Component[Selection]
	Move      Component[Move]
}

func (e Entity) Bounds() Rectangle {
	return Rectangle{
		e.Position.Value,
		e.Position.Value.Add(e.Image.Value.Bounds().Size()),
	}
}

func Draw(screen *ebiten.Image, e *Entity) {
	if e.Image.IsEnabled && e.Position.IsEnabled {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(e.Position.Value.X), float64(e.Position.Value.Y))
		screen.DrawImage(e.Image.Value, op)
	}
}
