package main

import "golang.org/x/exp/slog"

type PatrolOrder struct {
	origin      Point
	destination Point
}

func (o *PatrolOrder) Update(e *Entity, g *Game) {
	if e.Position.IsEnabled && e.Move.IsEnabled {
		if !e.Move.Value.IsActive {
			moveMap := g.getMoveMap()
			if e.Position.Value == o.destination {
				e.StartMove(o.origin, moveMap)
			} else {
				e.StartMove(o.destination, moveMap)
			}
		}
	}
}

func Patrol(e *Entity, destination Point) Order {
	if e.Selection.IsEnabled && e.Position.IsEnabled && e.Move.IsEnabled && e.Order.IsEnabled {
		if e.Selection.Value.IsSelected {
			origin := e.Position.Value
			slog.Info("patrolling between", slog.String("origin", origin.String()), slog.String("destination", destination.String()))
			e.Order.Value = &PatrolOrder{origin: origin, destination: destination}
		}
	}
	return &PatrolOrder{origin: e.Position.Value, destination: destination}
}
