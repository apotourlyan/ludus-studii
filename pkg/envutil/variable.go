package envutil

type Variable[T any] interface {
	Value() T
}
