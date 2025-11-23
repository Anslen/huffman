package main

type wordsFrequence = map[byte]int
type huffmanTree = Tree[charFrequence]
type HuffmanCodes = map[byte]HuffmanCode

type charFrequence struct {
	char      int16
	frequence int
}

type HuffmanCode struct {
	Code  uint64
	Width uint8
}

type huffmanTreeInfo struct {
	node *huffmanTree
	code HuffmanCode
}

// get frequence of each byte in the string
func getFrequence(str string) (ret wordsFrequence) {
	ret = make(wordsFrequence)
	data := []byte(str)
	for _, b := range data {
		ret[b] = ret[b] + 1
	}
	return ret
}

// build huffman tree from frequence map
func frequenceToTree(frequence wordsFrequence) (ret *huffmanTree) {
	// store [char, frequence]
	priority_queue := NewPriorityQueue[charFrequence](func(a1, a2 any) bool {
		return a1.(charFrequence).frequence < a2.(charFrequence).frequence
	})
	for key, value := range frequence {
		priority_queue.Push(charFrequence{int16(key), value})
	}

	// handle single char: create an internal root
	if priority_queue.Size() == 1 {
		pair, _ := priority_queue.Pop()
		parent := NewTree(charFrequence{int16(-1), pair.frequence})
		parent.Left = NewTree(pair)
		return parent
	}

	// map from frequence to tree
	trees := make(map[int16]*huffmanTree)
	treeIndex := int16(-1)
	for priority_queue.Size() > 1 {
		pairLeft, _ := priority_queue.Pop()
		pairRight, _ := priority_queue.Pop()

		treeLeft := trees[int16(pairLeft.char)]
		treeRight := trees[int16(pairRight.char)]
		// if not exist creat tree node
		if treeLeft == nil {
			treeLeft = NewTree(pairLeft)
		}
		if treeRight == nil {
			treeRight = NewTree(pairRight)
		}

		// store tree
		frequenceSum := pairLeft.frequence + pairRight.frequence
		parent := NewTree(charFrequence{treeIndex, frequenceSum})
		parent.Left = treeLeft
		parent.Right = treeRight
		trees[treeIndex] = parent
		treeIndex--

		priority_queue.Push(parent.Value)
	}

	// return first tree
	resultPair, _ := priority_queue.Top()
	ret = trees[resultPair.char]
	return ret
}

// convert huffman tree to huffman codes
func treeToCodes(tree *huffmanTree) (ret HuffmanCodes) {
	ret = make(HuffmanCodes)
	stack := NewStack[huffmanTreeInfo]()

	currentNode := tree
	currentCode := uint64(0)
	currentCodeWidth := uint8(0)

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
		if currentNode.Value.char >= 0 {
			ret[uint8(currentNode.Value.char)] = HuffmanCode{currentCode, currentCodeWidth}
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
	return ret
}

// build huffman tree from codes map
func codesToTree(codes HuffmanCodes) (ret *Tree[byte]) {
	ret = NewTree(byte(0))
	var current *Tree[byte] = ret
	var parent *Tree[byte]
	for char, code := range codes {
		reader := NewBitsReaderFromUint64(code.Code, int(code.Width))
		for i := 0; i < int(code.Width); i++ {
			parent = current
			bit, _ := reader.GetBit()
			if bit == 0 {
				current = current.Left
				if current == nil {
					current = NewTree(byte(0))
					parent.Left = current
				}
			} else {
				current = current.Right
				if current == nil {
					current = NewTree(byte(0))
					parent.Right = current
				}
			}
		}
		current.Value = char
		parent = nil
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
func HuffmanEncode(str string) (codes HuffmanCodes, codeWidth uint8, result []byte, resultWidth int64, ok bool) {
	if str == "" {
		return make(HuffmanCodes), 0, make([]byte, 0), 0, true
	}
	frequence := getFrequence(str)
	tree := frequenceToTree(frequence)
	codes = treeToCodes(tree)

	bitsRecorder := NewBitsRecorder()
	// iterate bytes (not runes) to match how frequence/count was computed
	data := []byte(str)
	for _, b := range data {
		huffman := codes[b]
		bitsRecorder.Add(uint64(huffman.Code), huffman.Width)
		resultWidth += int64(huffman.Width)
	}

	codeWidth = uint8(tree.Height() - 1)
	if codeWidth > 64 {
		return codes, codeWidth, result, resultWidth, false
	}
	if (codeWidth % 8) != 0 {
		codeWidth = (codeWidth/8 + 1) * 8
	}

	return codes, codeWidth, bitsRecorder.Result, resultWidth, true
}
