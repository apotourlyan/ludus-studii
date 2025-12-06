package secretutil

type Secret[T any] interface {
	Value() T
}
