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

func (tree *Tree[T]) Height() int {
	if tree == nil {
		return 0
	}
	leftHeight := tree.Left.Height()
	rightHeight := tree.Right.Height()
	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}
