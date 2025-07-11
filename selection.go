package main

import (
	"age_of_empires/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slog"
)

const SELECTION_HALO_WIDTH = 10

type Selection struct {
	IsSelected bool
	Halo       *ebiten.Image
}

func (e *Entity) SelectMultiple(cursor physics.Point, selection GlobalSelection) {
	if !e.Selection.IsEnabled {
		return
	}
	e.Selection.Value.IsSelected = false
	selectionBounds := physics.Rectangle{Min: selection.Start, Max: cursor}.Canon()
	if selectionBounds.Overlaps(e.Bounds()) {
		e.Selection.Value.IsSelected = true
		slog.Info("entity selected", slog.String("position", e.Position.Value.String()))
	}
}

func (e *Entity) SelectSingle(cursor physics.Point, canBeSelected bool) bool {
	if !e.Selection.IsEnabled {
		return false
	}
	e.Selection.Value.IsSelected = false
	if !canBeSelected || !cursor.In(e.Bounds()) {
		return false
	}
	e.Selection.Value.IsSelected = true
	slog.Info("entity selected", slog.String("position", e.Position.Value.String()))
	return true
}
