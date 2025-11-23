package main

import "fmt"

func ReadFile(str string) (result string, err error) {
	var bytes []byte = []byte(str)
	var reader *BitsReader = NewBitsReader(bytes, len(bytes)*8)

	// read huffman table
	var codes HuffmanCodes
	codes, err = readHuffmanTable(reader)
	if err != nil {
		return "", err
	}

	// read string
	return readString(reader, codesToTree(codes))
}

// read huffman table from reader
func readHuffmanTable(reader *BitsReader) (codes HuffmanCodes, err error) {
	codes = make(HuffmanCodes)
	var codeWidth uint8
	var char byte
	var code uint64

	for codeWidth > 0 {
		// ok flag for read
		var codeWidthOk, charOk, codeOk bool

		// read width, check read fail
		codeWidth, codeWidthOk = reader.GetUint8()
		if !codeWidthOk {
			return nil, fmt.Errorf("failed to read huffman table")
		}
		// check end of huffman table
		if codeWidth == 0 {
			break
		}

		// read char
		char, charOk = reader.GetByte()

		// read code
		// calculate and skip invalid zeros
		var codeStoreWidth = (codeWidth + 7) / 8 * 8
		var offset = codeStoreWidth - codeWidth
		reader.Seek(int(offset))
		code, codeOk = reader.GetNBits(int((codeWidth + 7) / 8 * 8))

		if !charOk || !codeOk {
			return nil, fmt.Errorf("failed to read huffman table")
		}

		// save code
		codes[char] = HuffmanCode{Code: code, Width: codeWidth}
	}
	return codes, nil
}

// read string from reader
func readString(reader *BitsReader, tree *Tree[byte]) (str string, err error) {
	// read data width
	var dataWidth uint64
	dataWidth, ok := reader.GetUint64()
	if !ok {
		return "", fmt.Errorf("failed to read data width")
	}

	var currentNode *Tree[byte] = tree
	var ret []byte = make([]byte, 0)
	for dataWidth > 0 {
		var bit uint8
		bit, ok = reader.GetBit()
		if !ok {
			return "", fmt.Errorf("failed to read data:\nno enough bits")
		}

		// move acrodding to bit
		if bit == 0 {
			currentNode = currentNode.Left
		} else {
			currentNode = currentNode.Right
		}

		// add data to ret when reach leaf node
		if currentNode.Left == nil && currentNode.Right == nil {
			ret = append(ret, currentNode.Value)
			currentNode = tree
			dataWidth--
		}

		if currentNode.Left == nil && currentNode.Right == nil {
			panic("Invalid huffman tree")
		}
	}

	// check reach end
	if currentNode != tree {
		return "", fmt.Errorf("invalid encoding data")
	}
	return string(ret), nil
}
