package main

// generic binary tree
type Tree[T any] struct {
	Value T
	Left  *Tree[T]
	Right *Tree[T]
}

// create a new tree node with given value
func NewTree[T any](value T) (ret *Tree[T]) {
	ret = new(Tree[T])
	ret.Value = value
	return ret
}

// get the height of the tree
func (tree *Tree[T]) Height() (height int) {
	if tree == nil {
		return 0
	}

	leftHeight := tree.Left.Height()
	rightHeight := tree.Right.Height()

	if leftHeight > rightHeight {
		return leftHeight + 1
	} else {
		return rightHeight + 1
	}
}
