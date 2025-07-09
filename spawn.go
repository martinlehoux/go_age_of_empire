package main

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"
	"log/slog"
	"time"
)

type spawnRequest struct {
	start time.Time
}

type Spawn struct {
	UnitResourceCost  int
	UnitSpawnDuration time.Duration
	Requests          []spawnRequest
	SpawnTarget       ecs.Component[physics.Point]
}

func NewSpawn(unitResourceCost int, unitSpawnDuration time.Duration) Spawn {
	return Spawn{
		UnitResourceCost:  unitResourceCost,
		UnitSpawnDuration: unitSpawnDuration,
		Requests:          make([]spawnRequest, 0),
	}
}

func (spawn *Spawn) AddRequest(g *Game) {
	if g.ResourceAmount < spawn.UnitResourceCost {
		return
	}
	g.ResourceAmount -= spawn.UnitResourceCost
	spawn.Requests = append(spawn.Requests, spawnRequest{start: time.Now()})
}

func UpdateSpawn(g *Game, spawn *ecs.Component[Spawn], position ecs.Component[physics.Point]) {
	if !spawn.IsEnabled || !position.IsEnabled {
		return
	}
	now := time.Now()
	if len(spawn.Value.Requests) == 0 || now.Sub(spawn.Value.Requests[0].start) < time.Duration(spawn.Value.UnitSpawnDuration) {
		return
	}
	spawnPosition, _ := g.Closest(position.Value, physics.AdjacentPoints(position.Value))
	unit := g.UnitBuilder.Build()
	unit.Position = ecs.C(spawnPosition)
	g.Entities = append(g.Entities, &unit)
	slog.Info("Spawned unit")
	if spawn.Value.SpawnTarget.IsEnabled {
		slog.Info("Unit has spawn target")
		unit.MainAction(g, spawn.Value.SpawnTarget.Value, g.entityAt(spawn.Value.SpawnTarget.Value), g.getMoveMap())
	}
	spawn.Value.Requests = spawn.Value.Requests[1:]
	if len(spawn.Value.Requests) > 0 {
		spawn.Value.Requests[0].start = now
	}
}
