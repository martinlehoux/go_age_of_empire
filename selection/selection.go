package selection

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const HALO_WIDTH = 10

type Priority int

const (
	Unit Priority = iota
	Building
)

type GlobalSelection struct {
	IsActive bool
	Start    physics.Point
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

func (m *Multiple) Preselect(position ecs.Component[physics.Point], bounds physics.Rectangle, selection *ecs.Component[Selection]) {
	if !selection.IsEnabled {
		return
	}
	selection.Value.IsSelected = false
	if m.Bounds.Overlaps(bounds) {
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

func (s *Single) Select(position ecs.Component[physics.Point], bounds physics.Rectangle, selection *ecs.Component[Selection]) {
	if !selection.IsEnabled {
		return
	}
	selection.Value.IsSelected = false
	if s.HasSelected || !s.Cursor.In(bounds) {
		return
	}
	selection.Value.IsSelected = true
}
