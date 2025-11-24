package main

import "fmt"

type huffmanTree = Tree[huffmanNode]
type HuffmanCodes = map[byte]HuffmanCode

// compare function for priority queue
func compareHuffmanTree(a1, a2 any) bool {
	if a1.(*huffmanTree).Value.frequence != a2.(*huffmanTree).Value.frequence {
		return a1.(*huffmanTree).Value.frequence < a2.(*huffmanTree).Value.frequence
	}
	if a1.(*huffmanTree).Value.char != a2.(*huffmanTree).Value.char {
		return a1.(*huffmanTree).Value.char < a2.(*huffmanTree).Value.char
	}
	return a1.(*huffmanTree).Value.index < a2.(*huffmanTree).Value.index
}

type huffmanNode struct {
	char      byte
	frequence int
	index     int
}

type HuffmanCode struct {
	Code  uint64
	Width uint8
}

// get frequence of each byte in the string
func getFrequence(str string) (ret map[byte]int) {
	ret = make(map[byte]int)
	for _, b := range []byte(str) {
		ret[b] = ret[b] + 1
	}
	return ret
}

// build huffman tree from frequence map
func frequenceToTree(frequence map[byte]int) (ret *huffmanTree) {
	// if single char, return tree with single node
	if len(frequence) == 1 {
		for char, frequence := range frequence {
			return NewTree(huffmanNode{char, frequence, 0})
		}
	}

	// priority queue to store trees
	var priority_queue *Priority_queue[*huffmanTree] = NewPriorityQueue[*huffmanTree](compareHuffmanTree)

	// add each char frequence to priority queue
	for char, frequence := range frequence {
		priority_queue.Push(NewTree(huffmanNode{char, frequence, 0}))
	}

	// node index
	var index int = 0
	// build tree
	for priority_queue.Size() > 1 {
		left, _ := priority_queue.Pop()
		right, _ := priority_queue.Pop()

		// create new internal node
		var parent *huffmanTree = NewTree(huffmanNode{char: 0, frequence: left.Value.frequence + right.Value.frequence, index: index})
		parent.Left = left
		parent.Right = right
		priority_queue.Push(parent)
		index++
	}

	// return first tree
	ret, _ = priority_queue.Pop()
	return ret
}

// convert huffman tree to huffman codes
//
// returns: HuffmanCodes map, code width in bits, success flag
//
// code width is align with byte (8 bits)
func treeToCodes(tree *huffmanTree) (ret HuffmanCodes, err error) {
	// store info for stack
	type huffmanTreeInfo struct {
		node *huffmanTree
		code HuffmanCode
	}

	if tree.Left == nil && tree.Right == nil {
		// single node tree
		// store byte to be 0 with width 1
		return HuffmanCodes{tree.Value.char: HuffmanCode{Code: 0, Width: 1}}, nil
	}

	ret = make(HuffmanCodes)
	var stack *Stack[huffmanTreeInfo] = NewStack[huffmanTreeInfo]()

	var currentNode *huffmanTree = tree
	var currentCode uint64 = uint64(0)
	var currentCodeWidth uint8 = uint8(0)

	for currentNode != nil || !stack.Empty() {
		// back
		if currentNode == nil {
			info, _ := stack.Pop()

			currentNode = info.node
			currentCode = info.code.Code
			currentCodeWidth = info.code.Width
			continue
		}

		// is leaf node: record huffman code
		if currentNode.Left == nil && currentNode.Right == nil {
			if currentCodeWidth > 64 {
				err = fmt.Errorf("huffman code longer than 64 bits")
				return nil, err
			}
			ret[currentNode.Value.char] = HuffmanCode{currentCode, currentCodeWidth}
			currentNode = nil
			continue
		}

		// calculate right child node
		childCode := (currentCode << 1) | 0x1
		stack.Push(huffmanTreeInfo{currentNode.Right, HuffmanCode{childCode, currentCodeWidth + 1}})

		// jump to left child node
		currentNode = currentNode.Left
		currentCode = currentCode << 1
		currentCodeWidth++
	}

	return ret, nil
}

func GetHuffmanCodes(str string) (codes HuffmanCodes, err error) {
	if str == "" {
		return make(HuffmanCodes), nil
	}
	frequence := getFrequence(str)
	tree := frequenceToTree(frequence)
	return treeToCodes(tree)
}

// build huffman tree without frequence from codes map
func GetHuffmanTree(codes HuffmanCodes) (ret *Tree[byte]) {
	ret = NewTree(byte(0))

	// insert each code to tree
	var current *Tree[byte] = ret
	for char, code := range codes {
		reader := NewBitsReaderFromUint64(code.Code, int(code.Width))
		for i := 0; i < int(code.Width); i++ {
			bit, _ := reader.GetBit()
			// if bit is 0, go left; else go right
			if bit == 0 {
				if current.Left == nil {
					current.Left = NewTree(byte(0))
				}
				current = current.Left
			} else {
				if current.Right == nil {
					current.Right = NewTree(byte(0))
				}
				current = current.Right
			}
		}
		current.Value = char
		current = ret
	}
	return ret
}
