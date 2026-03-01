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

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	cursor := physics.Point{X: x, Y: y}
	screen.Fill(soilColor)

	required := ecs.CM_Image | ecs.CM_Position
	for i, mask := range g.Mask {
		if mask&required == required {
			position := g.Position[i]
			image := g.Image[i]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(position.X), float64(position.Y))
			screen.DrawImage(image, op)
			if mask&ecs.CM_Selectable == ecs.CM_Selectable {
				selectable := g.Selectable[i]
				if selectable.IsSelected {
					opt := &ebiten.DrawImageOptions{}
					opt.GeoM.Translate(float64(position.X-selection.HaloWidth/2), float64(position.Y-selection.HaloWidth/2))
					screen.DrawImage(selectable.Halo, opt)
					if mask&ecs.CM_ResourceSource == ecs.CM_ResourceSource {
						source := g.ResourceSource[i]
						resourceText := fmt.Sprintf("%d", source.Remaining)
						op := &text.DrawOptions{}
						op.GeoM.Translate(float64(position.X+5), float64(position.Y+30))
						op.ColorScale.ScaleWithColor(color.White)
						text.Draw(screen, resourceText, &text.GoTextFace{
							Source: g.FaceSource,
							Size:   40,
						}, op)
					}
					if mask&ecs.CM_Spawn == ecs.CM_Spawn {
						spawn := g.Spawn[i]
						spawnText := fmt.Sprintf("%d", len(spawn.Requests))
						op := &text.DrawOptions{}
						op.GeoM.Translate(float64(position.X+5), float64(position.Y+70))
						op.ColorScale.ScaleWithColor(color.White)
						text.Draw(screen, spawnText, &text.GoTextFace{
							Source: g.FaceSource,
							Size:   40,
						}, op)
						if spawn.SpawnTarget.Ok {
							spawnTarget := spawn.SpawnTarget.Value
							vector.DrawFilledCircle(screen, float32(spawnTarget.X+50), float32(spawnTarget.Y+50), 30, Red, true)
						}
					}
					if mask&ecs.CM_Move == ecs.CM_Move {
						move := g.Move[i]
						dx := +image.Bounds().Dx() / 2
						dy := +image.Bounds().Dy() / 2
						last := position
						for _, point := range move.Path {
							vector.StrokeLine(screen, float32(last.X+dx), float32(last.Y+dy), float32(point.X+dx), float32(point.Y+dy), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
							last = point
						}
					}
				}
			}
		}
	}
	if g.GlobalSelection.IsActive(cursor) {
		g.GlobalSelection.Draw(screen, cursor)
	}

	drawTopBanner(screen, g)

	if len(g.GlobalSelection.Selected) > 0 {
		drawBottomBanner(screen, g)
	}
}

func drawTopBanner(screen *ebiten.Image, g *Game) {
	screenBounds := screen.Bounds()
	bannerHeight := float32(200)
	vector.DrawFilledRect(screen, 0, 0, float32(screenBounds.Dx()), bannerHeight, color.White, true)
	resourceText := fmt.Sprintf("Resources: %d", g.ResourceAmount)
	op := &text.DrawOptions{}
	op.GeoM.Translate(25, 25)
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, resourceText, &text.GoTextFace{
		Source: g.FaceSource,
		Size:   100,
	}, op)
}

func drawBottomBanner(screen *ebiten.Image, g *Game) {
	screenBounds := screen.Bounds()
	bannerHeight := float32(200)
	bannerTop := screenBounds.Dy() - int(bannerHeight)
	vector.DrawFilledRect(screen, 0, float32(bannerTop), float32(screenBounds.Dx()), bannerHeight, color.White, true)

	shortcuts := "[Right-click] Action  [A] Patrol  [Esc] Deselect"
	for _, i := range g.GlobalSelection.Selected {
		if g.Mask[i]&ecs.CM_Spawn == ecs.CM_Spawn {
			shortcuts += "  [S] Spawn unit"
			break
		}
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(25, float64(float32(bannerTop)+25))
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, shortcuts, &text.GoTextFace{
		Source: g.FaceSource,
		Size:   60,
	}, op)
}
