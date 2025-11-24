package main

import (
	"fmt"
	"os"
	"time"
)

type DecodeSize struct {
	Original int // in bytes
	Decoded  int // in bytes
}

func Decode(inputPath, outptuPath string) (decodeSize DecodeSize, decodeTime time.Duration, err error) {
	// record start time
	var startTime time.Time = time.Now()

	// read input file
	var bytes []byte
	bytes, err = os.ReadFile(inputPath)
	if err != nil {
		return decodeSize, decodeTime, fmt.Errorf("open input file %s failed:\n%v", inputPath, err.Error())
	}
	var reader *BitsReader = NewBitsReader(bytes, len(bytes)*8)

	// read huffman table
	var codes HuffmanCodes
	codes, err = readHuffmanTable(reader)
	if err != nil {
		return decodeSize, decodeTime, fmt.Errorf("read huffman table failed:\n%v", err.Error())
	}

	// build huffman tree and read string
	var tree *Tree[byte] = GetHuffmanTree(codes)
	var text []byte
	text, err = readString(reader, tree)

	// open output file
	var outputFile *os.File
	outputFile, err = OpenFile(outptuPath)
	if err != nil {
		return decodeSize, decodeTime, fmt.Errorf("open output file %s failed:\n%v", outptuPath, err.Error())
	}
	defer outputFile.Close()

	// write text to output file
	_, err = outputFile.Write(text)
	if err != nil {
		return decodeSize, decodeTime, fmt.Errorf("write decoded data to file %s failed:\n%v", outptuPath, err.Error())
	}

	// record size and time
	decodeSize = DecodeSize{
		Original: len(bytes),
		Decoded:  len(text),
	}
	decodeTime = time.Since(startTime)
	return decodeSize, decodeTime, nil
}

// read huffman table from reader
func readHuffmanTable(reader *BitsReader) (codes HuffmanCodes, err error) {
	codes = make(HuffmanCodes)
	var codeWidth uint8
	var char byte
	var code uint64

	for {
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
		code, codeOk = reader.GetNBits(int(codeWidth))

		if !charOk || !codeOk {
			return nil, fmt.Errorf("failed to read huffman table")
		}

		// save code
		codes[char] = HuffmanCode{Code: code, Width: codeWidth}
	}
	return codes, nil
}

// read string from reader
func readString(reader *BitsReader, tree *Tree[byte]) (text []byte, err error) {
	// read data width
	var dataWidth uint64
	dataWidth, ok := reader.GetUint64()
	if !ok {
		return nil, fmt.Errorf("failed to read data width")
	}

	var currentNode *Tree[byte] = tree
	text = make([]byte, 0)
	for i := uint64(0); i < dataWidth; i++ {
		var bit uint8
		bit, ok = reader.GetBit()
		if !ok {
			return nil, fmt.Errorf("failed to read data:\nno enough bits")
		}

		// move acrodding to bit
		if bit == 0 {
			currentNode = currentNode.Left
		} else {
			currentNode = currentNode.Right
		}

		// validate currentNode before accessing children
		if currentNode == nil {
			return nil, fmt.Errorf("invalid encoding data: reached nil node")
		}

		// add data to ret when reach leaf node
		if currentNode.Left == nil && currentNode.Right == nil {
			text = append(text, currentNode.Value)
			currentNode = tree
		}
	}

	// check reach end
	if currentNode != tree {
		return nil, fmt.Errorf("invalid encoding data")
	}
	return text, nil
}
