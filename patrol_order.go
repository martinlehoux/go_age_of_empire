package main

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"

	"golang.org/x/exp/slog"
)

type PatrolOrder struct {
	origin      physics.Point
	destination physics.Point
}

func (o *PatrolOrder) Update(e *Entity, g *Game) {
	if e.Position.IsEnabled && e.Move.IsEnabled {
		if !e.Move.Value.IsActive {
			moveMap := g.getMoveMap()
			if e.Position.Value == o.destination {
				physics.StartMove(&e.Move, e.Position, o.origin, moveMap)
			} else {
				physics.StartMove(&e.Move, e.Position, o.destination, moveMap)
			}
		}
	}
}

func Patrol(position ecs.Component[physics.Point], move ecs.Component[physics.Move], order *ecs.Component[Order], destination physics.Point) {
	if !position.IsEnabled || !move.IsEnabled || !order.IsEnabled {
		return
	}
	origin := position.Value
	slog.Info("patrolling between", slog.String("origin", origin.String()), slog.String("destination", destination.String()))
	order.Value = &PatrolOrder{origin: origin, destination: destination}
}
