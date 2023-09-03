package main

import (
	"math"
)

type Move struct {
	IsActive    bool
	Origin      Point
	Destination Point
	Nodes       []Point
}

func NewMove(origin Point, destination Point) Move {
	nodes := []Point{}
	current := Point{
		origin.X/100*100 + 50,
		origin.Y/100*100 + 50,
	}
	for current != destination {
		next := current
		if math.Abs(float64(destination.X-current.X)) > math.Abs(float64(destination.Y-current.Y)) {
			next.X = current.X + int(math.Copysign(100, float64(destination.X-current.X)))
		} else {
			next.Y = current.Y + int(math.Copysign(100, float64(destination.Y-current.Y)))
		}
		nodes = append(nodes, next)
		current = next
	}
	return Move{IsActive: true, Origin: origin, Destination: destination, Nodes: nodes}
}

func (m *Move) Update(position Point, speed int) Point {
	if !m.IsActive {
		return position
	}
	next := m.Nodes[0]
	if Distance(position, next) < float64(speed) {
		if next == m.Destination {
			m.IsActive = false
		} else {
			m.Nodes = m.Nodes[1:]
		}
		return next
	}
	remainingMove := next.Sub(position)
	delta := remainingMove.Mul(speed).Div(int(Length(remainingMove)))
	return position.Add(delta)
}
