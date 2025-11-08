package main

type Stack[T any] struct {
	data []T
}

const stackInitCapacity = 128

func NewStack[T any]() *Stack[T] {
	ret := new(Stack[T])
	ret.data = make([]T, 0, stackInitCapacity)
	return ret
}

func (stack *Stack[T]) Push(data T) {
	stack.data = append(stack.data, data)
}

func (stack *Stack[T]) Pop() (T, bool) {
	length := len(stack.data)
	if length == 0 {
		var zero T
		return zero, false
	}
	ret := stack.data[length-1]
	stack.data = stack.data[:length-1]
	return ret, true
}

func (stack *Stack[T]) Top() (T, bool) {
	if stack.Empty() {
		var zero T
		return zero, false
	}
	return stack.data[len(stack.data)-1], true
}

func (stack *Stack[T]) Size() int {
	return len(stack.data)
}

func (stack *Stack[T]) Empty() bool {
	return len(stack.data) == 0
}
