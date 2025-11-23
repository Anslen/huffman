package main

import "fmt"

type huffmanTree = Tree[charFrequence]
type HuffmanCodes = map[byte]HuffmanCode

// compare function for priority queue
func compareHuffmanTree(a1, a2 any) bool {
	return a1.(charFrequence).frequence < a2.(charFrequence).frequence
}

type charFrequence struct {
	char      byte
	frequence int
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
			return NewTree[charFrequence](charFrequence{char, frequence})
		}
	}

	// priority queue to store trees
	var priority_queue *Priority_queue[*huffmanTree] = NewPriorityQueue[*huffmanTree](compareHuffmanTree)

	// add each char frequence to priority queue
	for char, frequence := range frequence {
		priority_queue.Push(NewTree[charFrequence](charFrequence{char, frequence}))
	}

	// build tree
	for priority_queue.Size() > 1 {
		left, _ := priority_queue.Pop()
		right, _ := priority_queue.Pop()

		// create new internal node
		var parent *huffmanTree
		parent = NewTree[charFrequence](charFrequence{char: 0, frequence: left.Value.frequence + right.Value.frequence})
		parent.Left = left
		parent.Right = right
		priority_queue.Push(parent)
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
func treeToCodes(tree *huffmanTree) (ret HuffmanCodes, ok bool) {
	// store info for stack
	type huffmanTreeInfo struct {
		node *huffmanTree
		code HuffmanCode
	}

	if tree.Left == nil && tree.Right == nil {
		// single node tree
		// store byte to be 0 with width 1
		return HuffmanCodes{tree.Value.char: HuffmanCode{Code: 0, Width: 1}}, true
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
				fmt.Errorf("Error: huffman code longer than 64 bits\n")
				return nil, false
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

	return ret, true
}

func GetHuffmanCodes(str string) (codes HuffmanCodes, ok bool) {
	if str == "" {
		return make(HuffmanCodes), true
	}
	frequence := getFrequence(str)
	tree := frequenceToTree(frequence)
	return treeToCodes(tree)
}

// build huffman tree without frequence from codes map
func codesToTree(codes HuffmanCodes) (ret *Tree[byte]) {
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

// decode huffman encoded data
func HuffmanDecode(bin []byte) (ret string) {
	// Check if input is empty
	if len(bin) == 0 {
		return ""
	}

	reader := NewBitsReader(bin, len(bin)*8)
	if reader == nil {
		return ""
	}

	codeStoreWidth, ok := reader.GetUint8()
	if !ok {
		return ""
	}

	// Validate code width
	if codeStoreWidth == 0 || codeStoreWidth > 64 {
		return ""
	}

	// get huffman code
	codes := make(HuffmanCodes)
	char, ok := reader.GetByte()
	if !ok {
		return ""
	}

	code, ok := reader.GetNBits(int(codeStoreWidth))
	if !ok {
		return ""
	}

	codeWidth, ok := reader.GetUint8()
	if !ok {
		return ""
	}

	// Read code table with validation
	for codeWidth != 0 {
		codes[char] = HuffmanCode{code, codeWidth}

		char, ok = reader.GetByte()
		if !ok {
			return ""
		}

		code, ok = reader.GetNBits(int(codeStoreWidth))
		if !ok {
			return ""
		}

		codeWidth, ok = reader.GetUint8()
		if !ok {
			return ""
		}
	}

	// Validate that we have at least one code
	if len(codes) == 0 {
		return ""
	}

	// convert to huffman tree
	tree := codesToTree(codes)
	if tree == nil {
		return ""
	}

	result := make([]byte, 0)
	// find chars
	current := tree
	width, ok := reader.GetInt64()
	if !ok {
		return ""
	}

	// Validate width is not negative
	if width < 0 {
		return ""
	}

	for width != 0 {
		bit, ok := reader.GetBit()
		if !ok {
			return ""
		}

		width--
		if bit == 0 {
			current = current.Left
		} else {
			current = current.Right
		}

		// Check if we reached an invalid path in the tree
		if current == nil {
			return ""
		}

		// is leaf node: both children nil
		if current.Left == nil && current.Right == nil {
			result = append(result, current.Value)
			current = tree
		}
	}

	// return
	ret = string(result)
	return ret
}

// encode string using huffman coding
func HuffmanEncode(str string, codes HuffmanCodes) (result []byte, resultWidth int64) {
	if str == "" {
		return make([]byte, 0), 0
	}

	var bitsRecorder *BitsRecorder = NewBitsRecorder()
	// write each byte to bits recorder
	for _, char := range []byte(str) {
		var huffman HuffmanCode = codes[char]
		bitsRecorder.Add(huffman.Code, huffman.Width)
		resultWidth += int64(huffman.Width)
	}

	return bitsRecorder.Result, resultWidth
}
