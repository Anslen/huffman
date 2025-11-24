package main

import (
	"fmt"
	"io"
)

// write to output file
//
// format:
//
//	n group of:
//	    1 byte   : code width (in bits)
//	    1 byte   : character
//	    m bytes  : code, low bits valid, MSB first (m = code store width / 8, rounded up)
//	1 byte   : 0 (end of table)
//	8 bytes  : encoded data width (in bits)
//	n bytes  : encoded data
func WriteEncodeToFile(file io.Writer, str string, codes HuffmanCodes) (huffmanTableSize int, dataSize int, err error) {
	// write huffman table
	huffmanTableSize, err = writeHuffmanTable(file, codes)
	if err != nil {
		return huffmanTableSize, huffmanTableSize, err
	}

	// write string
	dataSize, err = writeString(file, str, codes)

	return huffmanTableSize, dataSize, err
}

// write huffman table to file
//
// return size written(in bytes) and ok
//
// format:
//
//	n group of:
//	    1 byte   : code width (in bits)
//	    1 byte   : character
//	    m bytes  : code, low bits valid, MSB first (m = code store width / 8, rounded up)
//	1 byte   : 0 (end of table)
func writeHuffmanTable(file io.Writer, codes HuffmanCodes) (size int, err error) {
	var recorder *BitsRecorder = NewBitsRecorder()

	// write each code
	for char, code := range codes {
		// code width
		recorder.Add(uint64(code.Width), 8)
		// character
		recorder.Add(uint64(char), 8)
		// code
		var codeStoreLength uint8 = (code.Width + 7) / 8 * 8
		recorder.Add(code.Code, codeStoreLength)
	}
	// end of table
	recorder.Add(0, 8)

	// write to file
	size, err = file.Write(recorder.Result())
	if err != nil {
		err = fmt.Errorf("write huffman table to file failed: %w", err)
		return size, err
	}
	return size, err
}

// write string to file using huffman coding
//
// return size written(in bytes) and ok
//
// format:
//
//	8 bytes  : encoded data width (in bits)
//	n bytes  : encoded data
func writeString(file io.Writer, str string, codes HuffmanCodes) (size int, err error) {
	// encode data
	var dataRecorder *BitsRecorder = NewBitsRecorder()
	// encoded data width (in bits)
	var dataWidth int64 = 0

	// write each byte to bits recorder
	for _, char := range []byte(str) {
		var huffman HuffmanCode = codes[char]
		dataRecorder.Add(huffman.Code, huffman.Width)
		dataWidth += int64(huffman.Width)
	}

	var widthRecorder *BitsRecorder = NewBitsRecorder()
	// write width info
	widthRecorder.Add(uint64(dataWidth), 64)

	// write to file
	size, err = file.Write(widthRecorder.Result())
	if err != nil {
		err = fmt.Errorf("write encoded data width to file failed:\n%w", err)
		return size, err
	}
	var widthWritenLength = size

	size, err = file.Write(dataRecorder.Result())
	if err != nil {
		err = fmt.Errorf("write encoded data to file failed:\n%w", err)
		return size + widthWritenLength, err
	}
	return size + widthWritenLength, err
}
