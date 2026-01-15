package main

import (
	"log/slog"
	"time"

	"age_of_empires/ecs"
	"age_of_empires/monad"
	"age_of_empires/physics"
)

type SpawnRequest struct {
	Start time.Time
}

type Spawn struct {
	UnitResourceCost  int
	UnitSpawnDuration time.Duration
	Requests          []SpawnRequest
	SpawnTarget       monad.Option[physics.Point]
}

func NewSpawn(unitResourceCost int, unitSpawnDuration time.Duration) Spawn {
	return Spawn{
		UnitResourceCost:  unitResourceCost,
		UnitSpawnDuration: unitSpawnDuration,
		Requests:          make([]SpawnRequest, 0),
	}
}

func UpdateSpawnSystem(g *Game, now time.Time, masks []ecs.Mask, spawns []Spawn, positions []physics.Point) {
	required := ecs.CM_Position | ecs.CM_Spawn
	for i, mask := range masks {
		if mask&required == required {
			spawn := &spawns[i]
			position := positions[i]
			if len(spawn.Requests) == 0 || now.Sub(spawn.Requests[0].Start) < time.Duration(spawn.UnitSpawnDuration) {
				continue
			}
			spawnPosition, _ := physics.Closest(position, physics.AdjacentPoints(position), g.getMoveMap())
			unit := g.Append(g.UnitBuilder.WithPosition(spawnPosition).Build())
			slog.Info("Spawned unit")
			if spawn.SpawnTarget.Ok {
				slog.Info("Unit has spawn target")
				MainAction(g, unit, spawn.SpawnTarget.Value, g.getMoveMap())
			}
			spawn.Requests = spawn.Requests[1:]
			if len(spawn.Requests) > 0 {
				spawn.Requests[0].Start = now
			}
		}
	}
}
