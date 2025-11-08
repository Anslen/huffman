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

func getFrequence(str string) (ret wordsFrequence, ok bool) {
	ret = make(wordsFrequence)
	data := []byte(str)
	for _, b := range data {
		ret[b] = ret[b] + 1
	}
	return ret, true
}

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

func codesToTree(codes HuffmanCodes) (ret *Tree[byte]) {
	ret = NewTree(byte(0))
	var current *Tree[byte] = ret
	var parent *Tree[byte]
	for char, code := range codes {
		reader := NewBitsReaderFromUint64(code.Code, 64)
		reader.Seek(int(64 - code.Width))
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

func HuffmanDecode(bin []byte) (ret string) {
	reader := NewBitsReader(bin, len(bin)*8)
	// get huffman code
	codes := make(HuffmanCodes)
	char, _ := reader.GetByte()
	code, _ := reader.GetUint64()
	codeWidth, _ := reader.GetUint8()
	for codeWidth != 0 {
		codes[char] = HuffmanCode{code, codeWidth}
		char, _ = reader.GetByte()
		code, _ = reader.GetUint64()
		codeWidth, _ = reader.GetUint8()
	}

	// convert to huffman tree
	tree := codesToTree(codes)

	result := make([]byte, 0)
	// find chars
	current := tree
	width, _ := reader.GetInt64()
	for width != 0 {
		bit, _ := reader.GetBit()
		width--
		if bit == 0 {
			current = current.Left
		} else {
			current = current.Right
		}

		// is leaf node save result
		if current.Value != 0 {
			result = append(result, current.Value)
			current = tree
		}
	}

	// return
	ret = string(result)
	return ret
}

func HuffmanEncode(str string) (codes HuffmanCodes, result []byte, width int64, ok bool) {
	frequence, ok := getFrequence(str)
	if !ok {
		return codes, result, width, false
	}
	tree := frequenceToTree(frequence)
	codes = treeToCodes(tree)

	bitsRecorder := NewBitsRecorder()
	for _, char := range str {
		huffman := codes[byte(char)]
		bitsRecorder.Add(uint64(huffman.Code), huffman.Width)
		width += int64(huffman.Width)
	}

	return codes, bitsRecorder.Result, width, true
}
