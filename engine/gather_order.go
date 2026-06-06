package engine

import (
	"math"
	"time"

	"age_of_empires/ecs"
	"age_of_empires/physics"

	"golang.org/x/exp/slog"
)

func UpdateGatherSystem(g *Game, now time.Time, moveMap physics.MoveMap, masks []ecs.Mask, orders []OrderKind, positions []physics.Point, moves []physics.Move, gatherers []ResourceGatherer) {
	required := ecs.CM_ResourceGatherer | ecs.CM_Position | ecs.CM_Move | ecs.CM_Order
	for i, mask := range masks {
		if mask&required == required {
			order := &orders[i]
			position := positions[i]
			move := &moves[i]
			gatherer := &gatherers[i]
			if *order == OrderKindGather {
				if gatherer.CurrentVolume == gatherer.MaxCapacity {
					updateGathererFull(g, getAllStorageDockings(g), moveMap, gatherer, order, move, position)
				} else if physics.Distance(position, g.Position[gatherer.Source]) <= 100 {
					updateGathering(now, gatherer, &g.ResourceSource[gatherer.Source])
				} else if !move.IsActive {
					startMoveToSource(moveMap, gatherer, move, order, position, g.Position[gatherer.Source])
				}
			}
		}
	}
}

func updateGathering(now time.Time, gatherer *ResourceGatherer, source *ResourceSource) {
	if gatherer.CurrentVolume >= gatherer.MaxCapacity || gatherer.LastPickupTime.Add(200*time.Millisecond).After(now) {
		return
	}
	gatherer.LastPickupTime = now
	gatherer.CurrentVolume += 1
	source.Remaining -= 1
}

func updateGathererFull(g *Game, storageDockings []physics.Point, moveMap physics.MoveMap, gatherer *ResourceGatherer, order *OrderKind, move *physics.Move, position physics.Point) {
	if move.IsActive {
		return
	}
	destination, distance := physics.Closest(position, storageDockings, moveMap)
	if distance == math.MaxInt {
		slog.Info("no accessible storage, canceling gather order")
		*order = OrderKindNone
		return
	}
	slog.Info("storage target", slog.String("storage", destination.String()), slog.Int("distance", distance))
	if distance > 1 {
		slog.Info("moving to storage", slog.String("storage", destination.String()))
		*move = physics.NewMove(position, destination, moveMap)
		return
	}
	g.ResourceAmount += gatherer.CurrentVolume
	gatherer.CurrentVolume = 0
	startMoveToSource(moveMap, gatherer, move, order, position, g.Position[gatherer.Source])
}

func startMoveToSource(moveMap physics.MoveMap, gatherer *ResourceGatherer, move *physics.Move, order *OrderKind, position physics.Point, sourcePosition physics.Point) {
	sourceDockings := physics.AdjacentPoints(sourcePosition)
	destination, distance := physics.Closest(position, sourceDockings, moveMap)
	if distance == math.MaxInt {
		slog.Info("no accessible storage, canceling gather order")
		*order = OrderKindNone
		return
	}
	*move = physics.NewMove(position, destination, moveMap)
}

func Gather(moveMap physics.MoveMap, position physics.Point, move *physics.Move, order *OrderKind, gatherer *ResourceGatherer, source ecs.Entity, sourcePosition physics.Point) {
	slog.Info("gathering from", slog.String("source", sourcePosition.String()))
	*order = OrderKindGather
	gatherer.Source = source
	startMoveToSource(moveMap, gatherer, move, order, position, sourcePosition)
}
