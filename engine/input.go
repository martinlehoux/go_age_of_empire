package engine

import (
	"age_of_empires/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	Cursor       physics.Point
	ScreenCursor physics.Point

	EscapePressed bool
	PatrolPressed bool

	LeftMouseDown bool
	LeftMouseUp   bool
	RightMouseUp  bool

	SpawnPressed bool
}

func CollectInput(camera *Camera) Input {
	sx, sy := ebiten.CursorPosition()
	return Input{
		Cursor:        camera.ScreenToWorld(sx, sy),
		ScreenCursor:  physics.Point{X: sx, Y: sy},
		EscapePressed: inpututil.IsKeyJustReleased(ebiten.KeyEscape),
		PatrolPressed: inpututil.IsKeyJustReleased(KeyPatrolRequest),
		LeftMouseDown: inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
		LeftMouseUp:   inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft),
		RightMouseUp:  inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight),
		SpawnPressed:  inpututil.IsKeyJustReleased(KeySpawnRequest),
	}
}
