package main

import (
	"math"

	"age_of_empires/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 640
	WorldWidth   = 3200
	WorldHeight  = 2400

	cameraZoomMin     = 0.25
	cameraZoomMax     = 3.0
	cameraZoomStep    = 0.10
	cameraScrollSpeed = 400.0
	cameraDrag        = 8.0
)

type Camera struct {
	X, Y float64
	Zoom float64

	vx, vy float64

	dragActive                         bool
	dragStartScreenX, dragStartScreenY int
	dragStartCamX, dragStartCamY       float64
}

func (c *Camera) WorldToScreen(p physics.Point) (float64, float64) {
	sx := (float64(p.X)-c.X)*c.Zoom + ScreenWidth/2
	sy := (float64(p.Y)-c.Y)*c.Zoom + ScreenHeight/2
	return sx, sy
}

func (c *Camera) ScreenToWorld(sx, sy int) physics.Point {
	zoom := c.Zoom
	if zoom == 0 {
		zoom = 1
	}
	wx := float64(sx-ScreenWidth/2)/zoom + c.X
	wy := float64(sy-ScreenHeight/2)/zoom + c.Y
	return physics.Point{X: int(math.Round(wx)), Y: int(math.Round(wy))}
}

func (c *Camera) GeoM() ebiten.GeoM {
	var m ebiten.GeoM
	m.Translate(-c.X, -c.Y)
	m.Scale(c.Zoom, c.Zoom)
	m.Translate(ScreenWidth/2, ScreenHeight/2)
	return m
}

func (c *Camera) clampToBounds() {
	halfW := (ScreenWidth / 2.0) / c.Zoom
	halfH := (ScreenHeight / 2.0) / c.Zoom

	minX, maxX := halfW, float64(WorldWidth)-halfW
	if minX > maxX {
		minX = WorldWidth / 2
		maxX = WorldWidth / 2
	}
	minY, maxY := halfH, float64(WorldHeight)-halfH
	if minY > maxY {
		minY = WorldHeight / 2
		maxY = WorldHeight / 2
	}

	if c.X < minX {
		c.X = minX
		if c.vx < 0 {
			c.vx = 0
		}
	}
	if c.X > maxX {
		c.X = maxX
		if c.vx > 0 {
			c.vx = 0
		}
	}
	if c.Y < minY {
		c.Y = minY
		if c.vy < 0 {
			c.vy = 0
		}
	}
	if c.Y > maxY {
		c.Y = maxY
		if c.vy > 0 {
			c.vy = 0
		}
	}
}

func (c *Camera) Update(dt float64) {
	zoom := c.Zoom
	if zoom == 0 {
		zoom = 1
	}

	// Scroll wheel zoom toward cursor
	_, wy := ebiten.Wheel()
	if wy != 0 {
		mx, my := ebiten.CursorPosition()
		// cursor world pos before zoom
		wx0 := float64(mx-ScreenWidth/2)/zoom + c.X
		wy0 := float64(my-ScreenHeight/2)/zoom + c.Y

		newZoom := math.Max(cameraZoomMin, math.Min(cameraZoomMax, zoom+wy*cameraZoomStep))
		c.Zoom = newZoom
		zoom = newZoom

		// translate so cursor stays fixed
		wx1 := float64(mx-ScreenWidth/2)/zoom + c.X
		wy1 := float64(my-ScreenHeight/2)/zoom + c.Y
		c.X -= wx1 - wx0
		c.Y -= wy1 - wy0
	}

	// Middle-mouse drag
	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		c.dragActive = true
		c.dragStartScreenX = mx
		c.dragStartScreenY = my
		c.dragStartCamX = c.X
		c.dragStartCamY = c.Y
	}
	if c.dragActive {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			c.X = c.dragStartCamX - float64(mx-c.dragStartScreenX)/zoom
			c.Y = c.dragStartCamY - float64(my-c.dragStartScreenY)/zoom
		} else {
			c.dragActive = false
		}
	}

	// Arrow key pan (velocity-based)
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy += 1
	}
	if dx != 0 && dy != 0 {
		inv := 1.0 / math.Sqrt2
		dx *= inv
		dy *= inv
	}
	c.vx += dx * cameraScrollSpeed
	c.vy += dy * cameraScrollSpeed

	// Apply velocity and exponential drag
	c.X += c.vx * dt
	c.Y += c.vy * dt
	decay := math.Exp(-cameraDrag * dt)
	c.vx *= decay
	c.vy *= decay

	c.clampToBounds()
}
