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
	if e.Selection.IsEnabled {
		e.Selection.Value.IsSelected = false
		selectionBounds := physics.Rectangle{Min: selection.Start, Max: cursor}.Canon()
		if selectionBounds.Overlaps(e.Bounds()) {
			e.Selection.Value.IsSelected = true
			slog.Info("entity selected", slog.String("position", e.Position.Value.String()))
		}
	}
}

func (e *Entity) SelectSingle(cursor physics.Point, canBeSelected bool) bool {
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

func DrawSelection(screen *ebiten.Image, e *Entity) {
	if e.Image.IsEnabled && e.Position.IsEnabled && e.Selection.IsEnabled {
		if e.Selection.Value.IsSelected {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(float64(e.Position.Value.X-SELECTION_HALO_WIDTH/2), float64(e.Position.Value.Y-SELECTION_HALO_WIDTH/2))
			screen.DrawImage(e.Selection.Value.Halo, opt)
		}
	}
}
