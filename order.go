package main

import "age_of_empires/physics"

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
	if !e.Selection.IsEnabled || !e.Selection.Value.IsSelected {
		return
	}
	if entityAtDestination != nil && entityAtDestination.ResourceSource.IsEnabled {
		Gather(e, entityAtDestination, g)
		return
	}
	physics.StartMove(&e.Move, e.Position, destination, moveMap)
}
