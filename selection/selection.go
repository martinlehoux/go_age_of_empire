package selection

import (
	"image/color"
	"math"

	"age_of_empires/ecs"
	"age_of_empires/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const HaloWidth = 10

type Priority int

const (
	Unit Priority = iota
	Building
)

type GlobalSelection struct {
	isActive bool
	start    physics.Point
}

func (s GlobalSelection) IsActive(cursor physics.Point) bool {
	return s.isActive && physics.Distance(s.start, cursor) > 10
}

func (s *GlobalSelection) Clear() {
	s.isActive = false
	s.start = physics.Point{}
}

func (s *GlobalSelection) Start(point physics.Point) {
	s.isActive = true
	s.start = point
}

func (s *GlobalSelection) Select(cursor physics.Point) Multiple {
	s.isActive = false
	return NewMultiple(s.start, cursor)
}

func (s *GlobalSelection) Draw(dst *ebiten.Image, cursor physics.Point) {
	vector.StrokeRect(dst, float32(s.start.X), float32(s.start.Y), float32(cursor.X-s.start.X), float32(cursor.Y-s.start.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
}

type Selection struct {
	IsSelected bool
	Halo       *ebiten.Image
	Priority   Priority
}

type Multiple struct {
	Bounds      physics.Rectangle
	Preselected []*ecs.Component[Selection]
}

func NewMultiple(start physics.Point, end physics.Point) Multiple {
	return Multiple{
		Bounds:      physics.Rectangle{Min: start, Max: end}.Canon(),
		Preselected: []*ecs.Component[Selection]{},
	}
}

func (m *Multiple) Preselect(position ecs.Component[physics.Point], relBounds ecs.Component[physics.Rectangle], selection *ecs.Component[Selection]) {
	if !selection.IsEnabled || !position.IsEnabled || !relBounds.IsEnabled {
		return
	}
	selection.Value.IsSelected = false
	if m.Bounds.Overlaps(physics.Translate(relBounds.Value, position.Value)) {
		m.Preselected = append(m.Preselected, selection)
	}
}

func (m *Multiple) Select() {
	if len(m.Preselected) == 0 {
		return
	}
	highestPriority := Priority(math.MaxInt)
	for _, selection := range m.Preselected {
		if selection.Value.Priority < highestPriority {
			highestPriority = selection.Value.Priority
		}
	}
	for _, selection := range m.Preselected {
		if selection.Value.Priority == highestPriority {
			selection.Value.IsSelected = true
		}
	}
}

type Single struct {
	Cursor      physics.Point
	HasSelected bool
}

func NewSingle(cursor physics.Point) Single {
	return Single{
		Cursor:      cursor,
		HasSelected: false,
	}
}

func (s *Single) Select(position ecs.Component[physics.Point], relBounds ecs.Component[physics.Rectangle], selection *ecs.Component[Selection]) {
	if !selection.IsEnabled || !position.IsEnabled || !relBounds.IsEnabled {
		return
	}
	selection.Value.IsSelected = false
	if s.HasSelected || !s.Cursor.In(physics.Translate(relBounds.Value, position.Value)) {
		return
	}
	selection.Value.IsSelected = true
}
