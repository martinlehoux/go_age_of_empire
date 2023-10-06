package main

import (
	"math"
	"time"

	"golang.org/x/exp/slog"
)

type GatherOrder struct {
	source *Entity
}

func (o *GatherOrder) Update(e *Entity, g *Game) {
	if e.Position.IsEnabled && e.Move.IsEnabled && e.ResourceGatherer.IsEnabled && e.ResourceGatherer.Value.CurrentTarget != nil {
		gatherer := &e.ResourceGatherer.Value
		source := gatherer.CurrentTarget.ResourceSource.Value
		if gatherer.CurrentVolume == gatherer.MaxCapacity && !e.Move.Value.IsActive {
			pointsNextToStorage := []Point{}
			for _, entity := range g.Entities {
				if entity.Position.IsEnabled && entity.ResourceStorage.IsEnabled {
					pointsNextToStorage = append(pointsNextToStorage, AdjacentPoints(entity.Position.Value)...)
				}
			}
			nextToStorage, distance := g.Closest(e.Position.Value, pointsNextToStorage)
			if distance == math.MaxInt {
				slog.Info("no accessible storage, canceling gather order")
				e.Order.Value = nil
				return
			}
			slog.Info("storage target", slog.String("storage", nextToStorage.String()), slog.Int("distance", distance))
			if distance > 1 {
				slog.Info("moving to storage", slog.String("storage", nextToStorage.String()))
				e.StartMove(nextToStorage, g.getMoveMap())
			} else {
				gatherer.CurrentVolume = 0
				// TODO: increment storage
				pointsNextToSource := AdjacentPoints(gatherer.CurrentTarget.Position.Value)
				nextToSource, distance := g.Closest(e.Position.Value, pointsNextToSource)
				if distance == math.MaxInt {
					slog.Info("no accessible storage, canceling gather order")
					e.Order.Value = nil
					return
				}
				e.StartMove(nextToSource, g.getMoveMap())
			}
		} else if Distance(e.Position.Value, gatherer.CurrentTarget.Position.Value) <= 100 {
			now := time.Now()
			if gatherer.CurrentVolume < gatherer.MaxCapacity && gatherer.LastPickupTime.Add(200*time.Millisecond).Before(now) {
				gatherer.LastPickupTime = now
				gatherer.CurrentVolume += 1
				source.Remaining -= 1
				slog.Info("gathered", slog.Int("current_volume", gatherer.CurrentVolume))
			}
		}
	}
}

func Gather(e *Entity, source *Entity, g *Game) Order {
	if e.Selection.IsEnabled && e.Position.IsEnabled && e.Move.IsEnabled && e.ResourceGatherer.IsEnabled && e.Order.IsEnabled {
		if e.Selection.Value.IsSelected {
			slog.Info("gathering from", slog.String("source", source.Position.Value.String()))
			e.Order.Value = &GatherOrder{source: source}
			e.ResourceGatherer.Value.CurrentTarget = source
			pointsNextToSource := AdjacentPoints(source.Position.Value)
			nextToSource, distance := g.Closest(e.Position.Value, pointsNextToSource)
			if distance == math.MaxInt {
				slog.Info("no accessible storage, canceling gather order")
				e.Order.Value = nil
				return nil
			}
			e.StartMove(nextToSource, g.getMoveMap())
		}
	}
	return &GatherOrder{source: source}
}
