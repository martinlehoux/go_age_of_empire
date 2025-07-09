package main

import (
	"age_of_empires/physics"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
		if e.Move.Value.IsActive {
			last := e.Position.Value
			dx := +e.Image.Value.Bounds().Dx() / 2
			dy := +e.Image.Value.Bounds().Dy() / 2
			for _, point := range e.Move.Value.Path {
				vector.StrokeLine(screen, float32(last.X+dx), float32(last.Y+dy), float32(point.X+dx), float32(point.Y+dy), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
				last = point
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	cursor := physics.Point{X: x, Y: y}
	screen.Fill(soilColor)
	for _, e := range g.Entities {
		e.Draw(g, screen)
	}
	if g.Selection.IsActive {
		vector.StrokeRect(screen, float32(g.Selection.Start.X), float32(g.Selection.Start.Y), float32(cursor.X-g.Selection.Start.X), float32(cursor.Y-g.Selection.Start.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
	}
	bannerHeight := float32(200)
	vector.DrawFilledRect(screen, 0, 0, float32(3200), bannerHeight, color.White, true)

	resourceText := fmt.Sprintf("Resources: %d", g.ResourceAmount)
	op := &text.DrawOptions{}
	op.GeoM.Translate(25, 25)
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, resourceText, &text.GoTextFace{
		Source: g.FaceSource,
		Size:   100,
	}, op)
}
