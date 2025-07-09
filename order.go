package main

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"
	"log/slog"
)

type Order interface {
	Update(e *Entity, g *Game)
}

func (e *Entity) UpdateOrder(g *Game) {
	if !e.Order.IsEnabled || e.Order.Value == nil {
		return
	}
	e.Order.Value.Update(e, g)
}

func (e *Entity) MainAction(g *Game, destination physics.Point, entityAtDestination *Entity, moveMap physics.MoveMap) {
	if entityAtDestination != nil && entityAtDestination.ResourceSource.IsEnabled && e.ResourceGatherer.IsEnabled {
		Gather(e, entityAtDestination, g)
		return
	}
	if e.Move.IsEnabled {
		physics.StartMove(&e.Move, e.Position, destination, moveMap)
		return
	}
	if e.Spawn.IsEnabled {
		slog.Info("Setting spawn target")
		e.Spawn.Value.SpawnTarget = ecs.C(destination)
		return
	}
}
