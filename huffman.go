package main

type wordsFrequence = map[rune]int
type pair = Pair[int64, int]
type huffmanTree = Tree[pair]
type HuffmanCode = Pair[uint64, int8]
type HuffmanCodes = map[rune]HuffmanCode

type huffmanTreeInfo struct {
	node      *huffmanTree
	code      uint64
	codeWidth int8
}

func getFrequence(str string) (ret wordsFrequence) {
	ret = make(wordsFrequence)
	for _, char := range str {
		frequence := ret[char]
		ret[char] = frequence + 1
	}
	return ret
}

func toHuffmanTree(frequence wordsFrequence) (ret *huffmanTree, ok bool) {
	// store [char, frequence]
	priority_queue := NewPriorityQueue[pair](func(a1, a2 any) bool {
		return a1.(pair).Second < a2.(pair).Second
	})
	for key, value := range frequence {
		priority_queue.Push(pair{int64(key), value})
	}

	// map from frequence to tree
	trees := make(map[int64]*huffmanTree)
	treeIndex := int64(-1)
	for priority_queue.Size() > 1 {
		pairLeft, _ := priority_queue.Pop()
		pairRight, _ := priority_queue.Pop()

		treeLeft := trees[pairLeft.First]
		treeRight := trees[pairRight.First]
		if treeLeft == nil {
			treeLeft = NewTree(pairLeft)
		}
		if treeRight == nil {
			treeRight = NewTree(pairRight)
		}

		frequenceSum := pairLeft.Second + pairRight.Second
		parent := NewTree(pair{treeIndex, frequenceSum})
		parent.Left = treeLeft
		parent.Right = treeRight
		trees[treeIndex] = parent
		treeIndex--

		priority_queue.Push(parent.Value)
	}

	if treeIndex >= 0 {
		return nil, false
	}

	resultPair, _ := priority_queue.Top()
	ret = trees[resultPair.First]
	return ret, true
}

func toHuffmanCode(tree *huffmanTree) (ret HuffmanCodes) {
	ret = make(HuffmanCodes)
	stack := NewStack[huffmanTreeInfo]()

	currentNode := tree
	currentCode := uint64(0)
	currentCodeWidth := int8(0)

	for currentNode != nil || !stack.Empty() {
		if currentNode == nil {
			info, _ := stack.Pop()

			currentNode = info.node
			currentCode = info.code
			currentCodeWidth = info.codeWidth
			continue
		}

		if currentNode.Value.First >= 0 {
			ret[rune(currentNode.Value.First)] = HuffmanCode{currentCode, currentCodeWidth}
			currentNode = nil
			continue
		}

		childCode := (currentCode << 1) | 0x1
		stack.Push(huffmanTreeInfo{currentNode.Right, childCode, currentCodeWidth + 1})

		currentNode = currentNode.Left
		currentCode = currentCode << 1
		currentCodeWidth++
	}
	return ret
}

func Huffman(str string) (codes HuffmanCodes, result []byte, width int64, ok bool) {
	frequence := getFrequence(str)
	tree, ok := toHuffmanTree(frequence)
	if !ok {
		return codes, result, width, false
	}
	codes = toHuffmanCode(tree)

	bitsRecorder := NewBitsRecorder()
	for _, char := range str {
		code := codes[char]
		bitsRecorder.Add(code.First, code.Second)
		width += int64(code.Second)
	}

	return codes, bitsRecorder.Result, width, true
}
