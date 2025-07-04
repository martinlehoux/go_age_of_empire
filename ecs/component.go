package ecs

type Component[T any] struct {
	IsEnabled bool
	Value     T
}

func C[T any](t T) Component[T] {
	return Component[T]{IsEnabled: true, Value: t}
}
