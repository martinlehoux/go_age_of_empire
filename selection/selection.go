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

type Selection struct {
	IsSelected bool
	Halo       *ebiten.Image
	Priority   Priority
}

type selectItem struct {
	ecs.Entity
	*Selection
}

type GlobalSelection struct {
	isActive bool
	start    physics.Point
	buffer   []selectItem
	Selected []ecs.Entity
}

func (s GlobalSelection) IsActive(cursor physics.Point) bool {
	return s.isActive
}

func (s GlobalSelection) IsArea(cursor physics.Point) bool {
	return physics.Distance(s.start, cursor) > 10
}

func (s *GlobalSelection) Clear() {
	s.isActive = false
	s.start = physics.Point{}
}

func (s *GlobalSelection) Start(point physics.Point) {
	s.isActive = true
	s.start = point
}

type SelectTarget struct {
	Selection *Selection
	Position  physics.Point
	Bounds    physics.Rectangle
}

func (s GlobalSelection) Bounds(cursor physics.Point) physics.Rectangle {
	return physics.Rectangle{
		Min: s.start,
		Max: cursor,
	}.Canon()
}

func (s *GlobalSelection) SelectMultiple(cursor physics.Point, masks []ecs.Mask, selections []Selection, positions []physics.Point, bounds []physics.Rectangle) {
	s.buffer = make([]selectItem, 0)
	s.Selected = make([]ecs.Entity, 0)
	required := ecs.CM_Position | ecs.CM_RelBounds | ecs.CM_Selection
	for i, mask := range masks {
		if mask&required == required {
			selection := &selections[i]
			position := positions[i]
			bound := bounds[i]

			selection.IsSelected = false
			if s.Bounds(cursor).Overlaps(physics.Translate(bound, position)) {
				s.buffer = append(s.buffer, selectItem{ecs.Entity(i), selection})
			}
		}
	}
	if len(s.buffer) == 0 {
		return
	}
	highestPriority := Priority(math.MaxInt)
	for _, selectable := range s.buffer {
		if selectable.Priority < highestPriority {
			highestPriority = selectable.Priority
		}
	}
	for _, selectable := range s.buffer {
		if selectable.Priority == highestPriority {
			selectable.IsSelected = true
			s.Selected = append(s.Selected, selectable.Entity)
		}
	}
	s.Clear()
}

func (s *GlobalSelection) SelectSingle(cursor physics.Point, masks []ecs.Mask, selections []Selection, positions []physics.Point, bounds []physics.Rectangle) {
	s.Selected = []ecs.Entity{}
	s.buffer = []selectItem{}
	required := ecs.CM_Position | ecs.CM_RelBounds | ecs.CM_Selection
	for i, mask := range masks {
		if mask&required == required {
			selection := &selections[i]
			position := positions[i]
			bound := bounds[i]
			selection.IsSelected = false
			if cursor.In(physics.Translate(bound, position)) {
				if len(s.buffer) == 0 || selection.Priority > s.buffer[0].Priority {
					s.buffer = []selectItem{{ecs.Entity(i), selection}}
				}
			}
		}
	}
	if len(s.buffer) > 0 {
		s.buffer[0].IsSelected = true
		s.Selected = []ecs.Entity{s.buffer[0].Entity}
	}
	s.Clear()
}

func (s *GlobalSelection) Unselect(selections []Selection) {
	for _, i := range s.Selected {
		selections[i].IsSelected = false
	}
	s.Selected = []ecs.Entity{}
}

func (s *GlobalSelection) Draw(dst *ebiten.Image, cursor physics.Point) {
	vector.StrokeRect(dst, float32(s.start.X), float32(s.start.Y), float32(cursor.X-s.start.X), float32(cursor.Y-s.start.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
}
