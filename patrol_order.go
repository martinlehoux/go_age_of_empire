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

func UpdatePatrolSystem(moveMap physics.MoveMap, masks []ecs.Mask, orders []OrderKind, patrols []PatrolOrder, positions []physics.Point, moves []physics.Move) {
	required := ecs.CM_Position | ecs.CM_Move | ecs.CM_Order | ecs.CM_PatrolOrder
	for i, mask := range masks {
		if mask&required == required {
			order := orders[i]
			patrol := patrols[i]
			position := positions[i]
			move := &moves[i]
			if order == OrderKindPatrol && !move.IsActive {
				if position == patrol.destination {
					*move = physics.NewMove(position, patrol.origin, moveMap)
				} else {
					*move = physics.NewMove(position, patrol.destination, moveMap)
				}
			}
		}
	}
}

func StartPatrolSystem(destination physics.Point, masks []ecs.Mask, positions []physics.Point, moves []physics.Move, orders []OrderKind, patrols []PatrolOrder) {
	required := ecs.CM_Position | ecs.CM_Order | ecs.CM_PatrolOrder
	for i, mask := range masks {
		if mask&required == required {
			position := positions[i]
			order := &orders[i]
			patrol := &patrols[i]
			slog.Info("patrolling between", slog.String("origin", position.String()), slog.String("destination", destination.String()))
			*patrol = PatrolOrder{origin: position, destination: destination}
			*order = OrderKindPatrol
		}
	}
}
