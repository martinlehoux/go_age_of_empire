package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slog"
)

type Selection struct {
	IsSelected    bool
	SelectedImage *ebiten.Image
}

func (e *Entity) SelectMultiple(cursor Point, selection GlobalSelection) {
	if e.Selection.IsEnabled {
		e.Selection.Value.IsSelected = false
		selectionBounds := Rectangle{selection.Start, cursor}.Canon()
		if selectionBounds.Overlaps(e.Bounds()) {
			e.Selection.Value.IsSelected = true
			slog.Info("entity selected", slog.String("position", e.Position.Value.String()))
		}
	}
}

func (e *Entity) SelectSingle(cursor Point, canBeSelected bool) bool {
	if e.Selection.IsEnabled {
		e.Selection.Value.IsSelected = false
		if canBeSelected && cursor.In(e.Bounds()) {
			e.Selection.Value.IsSelected = true
			slog.Info("entity selected", slog.String("position", e.Position.Value.String()))
			return true
		}
	}
	return false
}
