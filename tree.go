package main

type Tree[T any] struct {
	Value T
	Left  *Tree[T]
	Right *Tree[T]
}

func NewTree[T any](value T) (ret *Tree[T]) {
	ret = new(Tree[T])
	ret.Value = value
	return ret
}
