package main

type Compare func(any, any) bool

type priority_queue[T any] struct {
	data []T
	cmp  Compare
}

// create empty priority queue
//
// cmp: return true if a1 has higher priority than a2
//
// cmp(a1, a2) == true  => a1 higher priority than a2
func NewPriorityQueue[T any](cmp Compare) *priority_queue[T] {
	ret := new(priority_queue[T])
	ret.cmp = cmp
	return ret
}

// check if the priority queue is empty
func (queue *priority_queue[T]) Empty() bool {
	return len(queue.data) == 0
}

// get the size of the priority queue
func (queue *priority_queue[T]) Size() int {
	return len(queue.data)
}

// add a value to the priority queue
func (queue *priority_queue[T]) Push(value T) {
	holeIndex := len(queue.data)
	parent := (holeIndex - 1) / 2
	var empty T
	queue.data = append(queue.data, empty)
	for holeIndex != 0 && queue.cmp(value, queue.data[parent]) {
		queue.data[holeIndex] = queue.data[parent]
		holeIndex = parent
		parent = (holeIndex - 1) / 2
	}
	queue.data[holeIndex] = value
}

// get the top value of the priority queue
//
// return false if empty
func (queue *priority_queue[T]) Top() (retValue T, ok bool) {
	if len(queue.data) == 0 {
		return retValue, false
	}
	return queue.data[0], true
}

// get the index of the child to instead the hole
func insteadChild[T any](queue *priority_queue[T], index int) int {
	length := len(queue.data)
	left := index*2 + 1
	right := index*2 + 2
	if left >= length && right >= length {
		return -1
	}
	if left >= length {
		return right
	}
	if right >= length {
		return left
	}
	if queue.cmp(queue.data[left], queue.data[right]) {
		return left
	} else {
		return right
	}
}

// adjust the heap from holeIndex with value
func (queue *priority_queue[T]) adjustHeap(holeIndex int, value T) {
	insteadHoleIndex := insteadChild(queue, holeIndex)
	for insteadHoleIndex != -1 {
		queue.data[holeIndex] = queue.data[insteadHoleIndex]
		holeIndex = insteadHoleIndex
		insteadHoleIndex = insteadChild(queue, holeIndex)
	}
	// find position to put value
	parent := (holeIndex - 1) / 2
	for holeIndex != 0 && queue.cmp(value, queue.data[parent]) {
		queue.data[holeIndex] = queue.data[parent]
		holeIndex = parent
		parent = (holeIndex - 1) / 2
	}
	queue.data[holeIndex] = value
}

// remove and return the top value of the priority queue
//
// return false if empty
func (queue *priority_queue[T]) Pop() (retValue T, ok bool) {
	if len(queue.data) == 0 {
		return retValue, false
	}
	if len(queue.data) == 1 {
		retValue = queue.data[0]
		queue.data = queue.data[:0]
		return retValue, true
	}
	retValue = queue.data[0]
	value := queue.data[len(queue.data)-1]
	queue.data = queue.data[:len(queue.data)-1]
	queue.adjustHeap(0, value)
	return retValue, true
}
