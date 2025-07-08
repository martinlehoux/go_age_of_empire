package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (e *Entity) Draw(g *Game, screen *ebiten.Image) {
	if !e.Image.IsEnabled || !e.Position.IsEnabled {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(e.Position.Value.X), float64(e.Position.Value.Y))
	screen.DrawImage(e.Image.Value, op)
	if e.Selection.IsEnabled && e.Selection.Value.IsSelected {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(e.Position.Value.X-SELECTION_HALO_WIDTH/2), float64(e.Position.Value.Y-SELECTION_HALO_WIDTH/2))
		screen.DrawImage(e.Selection.Value.Halo, opt)
		if e.ResourceSource.IsEnabled {
			resourceText := fmt.Sprintf("%d", e.ResourceSource.Value.Remaining)
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(e.Position.Value.X+5), float64(e.Position.Value.Y+30))
			op.ColorScale.ScaleWithColor(color.White)
			text.Draw(screen, resourceText, &text.GoTextFace{
				Source: g.FaceSource,
				Size:   40,
			}, op)
		}
		if e.Spawn.IsEnabled {
			spawnText := fmt.Sprintf("%d", len(e.Spawn.Value.Requests))
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(e.Position.Value.X+5), float64(e.Position.Value.Y+70))
			op.ColorScale.ScaleWithColor(color.White)
			text.Draw(screen, spawnText, &text.GoTextFace{
				Source: g.FaceSource,
				Size:   40,
			}, op)
		}
	}
}
