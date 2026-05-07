package main

import (
	"fmt"
	"image/color"

	"age_of_empires/ecs"
	"age_of_empires/physics"
	"age_of_empires/selection"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) font(size float64) *text.GoTextFace {
	f := &text.GoTextFace{Source: g.FaceSource, Size: size}
	f.SetVariation(text.MustParseTag("wght"), 500)
	return f
}

func (g *Game) Draw(screen *ebiten.Image) {
	sx, sy := ebiten.CursorPosition()
	cursor := g.Camera.ScreenToWorld(sx, sy)
	screenCursor := physics.Point{X: sx, Y: sy}
	screen.Fill(soilColor)

	camGeoM := g.Camera.GeoM()

	required := ecs.CM_Image | ecs.CM_Position
	for i, mask := range g.Mask {
		if mask&required == required {
			position := g.Position[i]
			image := g.Image[i]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(position.X), float64(position.Y))
			op.GeoM.Concat(camGeoM)
			screen.DrawImage(image, op)
			if mask&ecs.CM_Selectable == ecs.CM_Selectable {
				selectable := g.Selectable[i]
				if selectable.IsSelected {
					opt := &ebiten.DrawImageOptions{}
					opt.GeoM.Translate(float64(position.X-selection.HaloWidth/2), float64(position.Y-selection.HaloWidth/2))
					opt.GeoM.Concat(camGeoM)
					screen.DrawImage(selectable.Halo, opt)
					if mask&ecs.CM_ResourceSource == ecs.CM_ResourceSource {
						source := g.ResourceSource[i]
						resourceText := fmt.Sprintf("%d", source.Remaining)
						op := &text.DrawOptions{}
						op.GeoM.Translate(float64(position.X+5), float64(position.Y+30))
						op.GeoM.Concat(camGeoM)
						op.ColorScale.ScaleWithColor(color.White)
						text.Draw(screen, resourceText, g.font(40), op)
					}
					if mask&ecs.CM_Spawn == ecs.CM_Spawn {
						spawn := g.Spawn[i]
						spawnText := fmt.Sprintf("%d", len(spawn.Requests))
						op := &text.DrawOptions{}
						op.GeoM.Translate(float64(position.X+5), float64(position.Y+70))
						op.GeoM.Concat(camGeoM)
						op.ColorScale.ScaleWithColor(color.White)
						text.Draw(screen, spawnText, g.font(40), op)
						if spawn.SpawnTarget.Ok {
							spawnTarget := spawn.SpawnTarget.Value
							ssx, ssy := g.Camera.WorldToScreen(spawnTarget)
							offset := 50 * float32(g.Camera.Zoom)
							vector.DrawFilledCircle(screen, float32(ssx)+offset, float32(ssy)+offset, 30*float32(g.Camera.Zoom), Red, true)
						}
					}
					if mask&ecs.CM_Move == ecs.CM_Move {
						move := g.Move[i]
						dx := image.Bounds().Dx() / 2
						dy := image.Bounds().Dy() / 2
						last := position
						strokeW := 10.0 * float32(g.Camera.Zoom)
						for _, point := range move.Path {
							lsx, lsy := g.Camera.WorldToScreen(physics.Point{X: last.X + dx, Y: last.Y + dy})
							psx, psy := g.Camera.WorldToScreen(physics.Point{X: point.X + dx, Y: point.Y + dy})
							vector.StrokeLine(screen, float32(lsx), float32(lsy), float32(psx), float32(psy), strokeW, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
							last = point
						}
					}
				}
			}
		}
	}
	if g.GlobalSelection.IsActive(cursor) {
		g.GlobalSelection.Draw(screen, screenCursor)
	}

	drawTopBanner(screen, g)

	if len(g.GlobalSelection.Selected) > 0 {
		drawBottomBanner(screen, g)
	}
}

func drawTopBanner(screen *ebiten.Image, g *Game) {
	screenBounds := screen.Bounds()
	bannerHeight := float32(50)
	vector.DrawFilledRect(screen, 0, 0, float32(screenBounds.Dx()), bannerHeight, color.White, true)
	resourceText := fmt.Sprintf("Resources: %d", g.ResourceAmount)
	op := &text.DrawOptions{}
	op.GeoM.Translate(8, 6)
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, resourceText, g.font(30), op)
}

func drawBottomBanner(screen *ebiten.Image, g *Game) {
	screenBounds := screen.Bounds()
	bannerHeight := float32(40)
	bannerTop := screenBounds.Dy() - int(bannerHeight)
	vector.DrawFilledRect(screen, 0, float32(bannerTop), float32(screenBounds.Dx()), bannerHeight, color.White, true)

	shortcuts := "[Right-click] Action  [Q] Patrol  [Esc] Deselect"
	for _, i := range g.GlobalSelection.Selected {
		if g.Mask[i]&ecs.CM_Spawn == ecs.CM_Spawn {
			shortcuts += "  [S] Spawn unit"
			break
		}
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(8, float64(float32(bannerTop)+6))
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, shortcuts, g.font(22), op)
}
