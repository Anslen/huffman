package main

type Stack[T any] struct {
	data []T
}

const stackInitCapacity = 128

// create empty stack
func NewStack[T any]() *Stack[T] {
	ret := new(Stack[T])
	ret.data = make([]T, 0, stackInitCapacity)
	return ret
}

// push data onto stack
func (stack *Stack[T]) Push(data T) {
	stack.data = append(stack.data, data)
}

// pop data from stack
//
// return false if stack is empty
func (stack *Stack[T]) Pop() (ret T, ok bool) {
	length := len(stack.data)
	if length == 0 {
		return ret, false
	}
	ret = stack.data[length-1]
	stack.data = stack.data[:length-1]
	return ret, true
}

// get the top data of stack
func (stack *Stack[T]) Top() (T, bool) {
	if stack.Empty() {
		var zero T
		return zero, false
	}
	return stack.data[len(stack.data)-1], true
}

// get the size of the stack
func (stack *Stack[T]) Size() int {
	return len(stack.data)
}

// check if the stack is empty
func (stack *Stack[T]) Empty() bool {
	return len(stack.data) == 0
}
