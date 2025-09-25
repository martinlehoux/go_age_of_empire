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

type GlobalSelection struct {
	isActive  bool
	start     physics.Point
	buffer    []*ecs.Component[Selection]
	selection []*ecs.Component[Selection]
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
	Selection *ecs.Component[Selection]
	Position  ecs.Component[physics.Point]
	Bounds    ecs.Component[physics.Rectangle]
}

func (s GlobalSelection) Bounds(cursor physics.Point) physics.Rectangle {
	return physics.Rectangle{
		Min: s.start,
		Max: cursor,
	}.Canon()
}

func (s *GlobalSelection) SelectMultiple(cursor physics.Point, components []SelectTarget) {
	s.buffer = make([]*ecs.Component[Selection], 0)
	s.selection = make([]*ecs.Component[Selection], 0)
	for _, target := range components {
		if target.Selection.IsEnabled && target.Position.IsEnabled && target.Bounds.IsEnabled {
			target.Selection.Value.IsSelected = false
			if s.Bounds(cursor).Overlaps(physics.Translate(target.Bounds.Value, target.Position.Value)) {
				s.buffer = append(s.buffer, target.Selection)
			}
		}
	}
	if len(s.buffer) == 0 {
		return
	}
	highestPriority := Priority(math.MaxInt)
	for _, selection := range s.buffer {
		if selection.Value.Priority < highestPriority {
			highestPriority = selection.Value.Priority
		}
	}
	for _, selection := range s.buffer {
		if selection.Value.Priority == highestPriority {
			selection.Value.IsSelected = true
			s.selection = append(s.selection, selection)
		}
	}
	s.Clear()
}

func (s *GlobalSelection) SelectSingle(cursor physics.Point, components []SelectTarget) {
	s.selection = []*ecs.Component[Selection]{}
	for _, target := range components {
		if target.Selection.IsEnabled && target.Position.IsEnabled && target.Bounds.IsEnabled {
			target.Selection.Value.IsSelected = false
			if cursor.In(physics.Translate(target.Bounds.Value, target.Position.Value)) {
				if len(s.selection) == 0 || target.Selection.Value.Priority > s.selection[0].Value.Priority {
					s.selection = []*ecs.Component[Selection]{target.Selection}
				}
			}
		}
	}
	s.Clear()
}

func (s *GlobalSelection) Unselect() {
	for _, selection := range s.selection {
		selection.Value.IsSelected = false
	}
	s.selection = []*ecs.Component[Selection]{}
}

func (s *GlobalSelection) Draw(dst *ebiten.Image, cursor physics.Point) {
	vector.StrokeRect(dst, float32(s.start.X), float32(s.start.Y), float32(cursor.X-s.start.X), float32(cursor.Y-s.start.Y), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
}
