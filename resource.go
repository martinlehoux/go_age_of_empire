package main

import "time"

type ResourceSource struct {
	Remaining int
}

type ResourceGatherer struct {
	MaxCapacity    int
	CurrentVolume  int
	CurrentTarget  *Entity
	LastPickupTime time.Time
}

type ResourceStorage struct{}
