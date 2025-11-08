package main

type wordsFrequence = map[byte]int
type huffmanTree = Tree[charFrequence]
type HuffmanCodes = map[byte]HuffmanCode

type charFrequence struct {
	char      int8
	frequence int
}

type HuffmanCode struct {
	code  uint8
	width int8
}

type huffmanTreeInfo struct {
	node *huffmanTree
	code HuffmanCode
}

func getFrequence(str string) (ret wordsFrequence, ok bool) {
	ret = make(wordsFrequence)
	for _, char := range str {
		if char > rune(127) {
			return make(wordsFrequence), false
		}
		frequence := ret[byte(char)]
		ret[byte(char)] = frequence + 1
	}
	return ret, true
}

func frequenceToTree(frequence wordsFrequence) (ret *huffmanTree) {
	// store [char, frequence]
	priority_queue := NewPriorityQueue[charFrequence](func(a1, a2 any) bool {
		return a1.(charFrequence).frequence < a2.(charFrequence).frequence
	})
	for key, value := range frequence {
		priority_queue.Push(charFrequence{int8(key), value})
	}

	// map from frequence to tree
	trees := make(map[int8]*huffmanTree)
	treeIndex := int8(-1)
	for priority_queue.Size() > 1 {
		pairLeft, _ := priority_queue.Pop()
		pairRight, _ := priority_queue.Pop()

		treeLeft := trees[pairLeft.char]
		treeRight := trees[pairRight.char]
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
	currentCode := uint8(0)
	currentCodeWidth := int8(0)

	for currentNode != nil || !stack.Empty() {
		// back
		if currentNode == nil {
			info, _ := stack.Pop()

			currentNode = info.node
			currentCode = info.code.code
			currentCodeWidth = info.code.width
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

func Huffman(str string) (codes HuffmanCodes, result []byte, width int64, ok bool) {
	frequence, ok := getFrequence(str)
	if !ok {
		return codes, result, width, false
	}
	tree := frequenceToTree(frequence)
	codes = treeToCodes(tree)

	bitsRecorder := NewBitsRecorder()
	for _, char := range str {
		huffman := codes[byte(char)]
		bitsRecorder.Add(uint64(huffman.code), huffman.width)
		width += int64(huffman.width)
	}

	return codes, bitsRecorder.Result, width, true
}
