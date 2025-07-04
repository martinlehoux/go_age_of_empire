package main

import (
	"age_of_empires/physics"
	"math"
	"time"

	"golang.org/x/exp/slog"
)

type GatherOrder struct {
	source *Entity
}

func (o *GatherOrder) Update(e *Entity, g *Game) {
	if !e.Position.IsEnabled || !e.Move.IsEnabled || !e.ResourceGatherer.IsEnabled || e.ResourceGatherer.Value.CurrentTarget == nil {
		return
	}

	gatherer := &e.ResourceGatherer.Value
	source := gatherer.CurrentTarget.ResourceSource.Value
	if gatherer.CurrentVolume == gatherer.MaxCapacity {
		o.updateGathererFull(e, g, gatherer)
		return
	}
	if physics.Distance(e.Position.Value, gatherer.CurrentTarget.Position.Value) <= 100 {
			o.updateGathering(gatherer, source)
			return
	}
	if !e.Move.Value.IsActive {
		o.startMoveToSource(e, gatherer, g)
	}
}

func (*GatherOrder) updateGathering(gatherer *ResourceGatherer, source ResourceSource) {
	now := time.Now()
	if gatherer.CurrentVolume >= gatherer.MaxCapacity || gatherer.LastPickupTime.Add(200*time.Millisecond).After(now) {
		return
	}
	gatherer.LastPickupTime = now
	gatherer.CurrentVolume += 1
	source.Remaining -= 1
	slog.Info("gathered", slog.Int("current_volume", gatherer.CurrentVolume))
}

func (o *GatherOrder) updateGathererFull(e *Entity, g *Game, gatherer *ResourceGatherer) {
	if e.Move.Value.IsActive {
		return
	}
	storageDockings := getAllStorageDockings(g)
	destination, distance := g.Closest(e.Position.Value, storageDockings)
	if distance == math.MaxInt {
		slog.Info("no accessible storage, canceling gather order")
		e.Order.Value = nil
		return
	}
	slog.Info("storage target", slog.String("storage", destination.String()), slog.Int("distance", distance))
	if distance > 1 {
		slog.Info("moving to storage", slog.String("storage", destination.String()))
		physics.StartMove(&e.Move,e.Position, destination, g.getMoveMap())
		return
	}
	gatherer.CurrentVolume = 0
	// TODO: increment storage
	o.startMoveToSource(e, gatherer, g)
}

func (o *GatherOrder) startMoveToSource(e *Entity, gatherer *ResourceGatherer, g *Game) {
	sourceDockings := physics.AdjacentPoints(gatherer.CurrentTarget.Position.Value)
	destination, distance := g.Closest(e.Position.Value, sourceDockings)
	if distance == math.MaxInt {
		slog.Info("no accessible storage, canceling gather order")
		e.Order.Value = nil
		return
	}
	physics.StartMove(&e.Move,e.Position, destination, g.getMoveMap())
}

func getAllStorageDockings(g *Game) []physics.Point {
	storageDockings := []physics.Point{}
	for _, entity := range g.Entities {
		if entity.Position.IsEnabled && entity.ResourceStorage.IsEnabled {
			storageDockings = append(storageDockings, physics.AdjacentPoints(entity.Position.Value)...)
		}
	}
	return storageDockings
}

func Gather(e *Entity, source *Entity, g *Game) Order {
	if e.Selection.IsEnabled && e.Position.IsEnabled && e.Move.IsEnabled && e.ResourceGatherer.IsEnabled && e.Order.IsEnabled {
		if e.Selection.Value.IsSelected {
			slog.Info("gathering from", slog.String("source", source.Position.Value.String()))
			e.Order.Value = &GatherOrder{source: source}
			e.ResourceGatherer.Value.CurrentTarget = source
			sourceDockings := physics.AdjacentPoints(source.Position.Value)
			destination, distance := g.Closest(e.Position.Value, sourceDockings)
			if distance == math.MaxInt {
				slog.Info("no accessible storage, canceling gather order")
				e.Order.Value = nil
				return nil
			}
			physics.StartMove(&e.Move,e.Position, destination, g.getMoveMap())
		}
	}
	return &GatherOrder{source: source}
}
