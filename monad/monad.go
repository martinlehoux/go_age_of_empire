package monad

type Option[T any] struct {
	Ok    bool
	Value T
}

func Some[T any](value T) Option[T] {
	return Option[T]{Ok: true, Value: value}
}
