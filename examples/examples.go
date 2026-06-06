package examples

import (
	"age_of_empires/engine"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type BuilderDeps struct {
	FaceSource  *text.GoTextFaceSource
	UnitBuilder engine.EntityBuilder
}

type Entry struct {
	Name        string
	Description string
	Build       func(*engine.Game, BuilderDeps)
}

var Registry = map[string]Entry{}
