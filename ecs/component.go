package ecs

type Mask uint16

const (
	CM_Position Mask = 1 << iota
	CM_RelBounds
	CM_Image
	CM_Selection
	CM_Move
	CM_Order
	CM_PatrolOrder
	CM_GatherOrder
	CM_ResourceGatherer
	CM_ResourceSource
	CM_ResourceStorage
	CM_Spawn
)
