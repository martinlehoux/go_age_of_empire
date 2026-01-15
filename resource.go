package main

import (
	"time"

	"age_of_empires/ecs"
)

type ResourceSource struct {
	Remaining int
}

type ResourceGatherer struct {
	MaxCapacity    int
	CurrentVolume  int
	Source         ecs.Entity
	LastPickupTime time.Time
}

type ResourceStorage struct{}
