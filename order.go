package main

import (
	"log/slog"

	"age_of_empires/ecs"
	"age_of_empires/monad"
	"age_of_empires/physics"
)

type OrderKind int

const (
	OrderKindNone OrderKind = iota
	OrderKindPatrol
	OrderKindGather
)

func MainAction(g *Game, mainEntity ecs.Entity, target physics.Point, moveMap physics.MoveMap) {
	mask := g.Mask[mainEntity]
	targetEntity := g.entityAt(target)
	gatherTargetRequired := ecs.CM_ResourceSource
	mainGatherRequired := ecs.CM_ResourceGatherer | ecs.CM_Move
	if targetEntity >= 0 && (g.Mask[targetEntity]&gatherTargetRequired == gatherTargetRequired) && (mask&mainGatherRequired == mainGatherRequired) {
		Gather(g.getMoveMap(), g.Position[mainEntity], &g.Move[mainEntity], &g.Order[mainEntity], &g.ResourceGatherer[mainEntity], targetEntity, g.Position[targetEntity])
		return
	}
	moveRequired := ecs.CM_Move | ecs.CM_Position
	if mask&moveRequired == moveRequired {
		position := g.Position[mainEntity]
		g.Order[mainEntity] = OrderKindNone
		g.Move[mainEntity] = physics.NewMove(position, target, moveMap)
		return
	}
	if mask&ecs.CM_Spawn == ecs.CM_Spawn {
		slog.Info("Setting spawn target")
		spawn := &g.Spawn[mainEntity]
		spawn.SpawnTarget = monad.Some(target)
		return
	}
}
