package main

import (
	"age_of_empires/ecs"
	"age_of_empires/physics"
	"time"
)

type spawnRequest struct {
	start time.Time
}

type Spawn struct {
	UnitResourceCost int
	UnitSpawnDuration time.Duration
	Requests []spawnRequest
}

func NewSpawn(unitResourceCost int, unitSpawnDuration time.Duration) Spawn {
	return Spawn{
		UnitResourceCost: unitResourceCost,
		UnitSpawnDuration: unitSpawnDuration,
		Requests: make([]spawnRequest, 0),
	}
}

func (spawn *Spawn) AddRequest(g *Game) {
	if g.ResourceAmount < spawn.UnitResourceCost {
		return
	}
	g.ResourceAmount -= spawn.UnitResourceCost
	spawn.Requests = append(spawn.Requests, spawnRequest{start: time.Now()})
}

func UpdateSpawn(g *Game,spawn *ecs.Component[Spawn], position ecs.Component[physics.Point]) {
	if !spawn.IsEnabled || !position.IsEnabled {
		return
	}
	now := time.Now()
	if len(spawn.Value.Requests) == 0 || now.Sub(spawn.Value.Requests[0].start) < time.Duration(spawn.Value.UnitSpawnDuration) {
		return
	}
	// TODO: builder
	var order Order
	spawnPosition, _ := g.Closest(position.Value, physics.AdjacentPoints(position.Value))
	unit := Entity{
		Position:         ecs.C(spawnPosition),
		Image:            ecs.C(g.personImage),
		Selection:        ecs.C(Selection{IsSelected: false, Halo: g.personSelectionHalo}),
		Move:             ecs.C(physics.Move{IsActive: false}),
		Order:            ecs.C(order),
		ResourceGatherer: ecs.C(ResourceGatherer{MaxCapacity: 15}),
	}
	g.Entities = append(g.Entities, &unit)
	spawn.Value.Requests = spawn.Value.Requests[1:]
	if len(spawn.Value.Requests) > 0 {
		spawn.Value.Requests[0].start = now
	}
}
