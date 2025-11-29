package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type EncodeSize struct {
	orininal     int // in bytes
	HuffmanTable int // in bytes
	EncodedData  int // in bytes
}

type EncodeTime struct {
	CodeGenTime   time.Duration // in milliseconds
	WriteFileTime time.Duration // in milliseconds
}

type BatchError struct {
	Path string
	Err  error
}

type BatchEncodeResult struct {
	InputPath    string
	OutputPath   string
	TotalCount   int
	SuccessCount int
	Time         time.Duration
	Errors       []BatchError
}

// write to output file
//
// return input size and output size(in bytes) and ok
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
func Encode(inputPath, outputPath string) (encodeSize EncodeSize, encodeTime EncodeTime, err error) {
	// record start time
	var startTime time.Time = time.Now()

	// read input file
	var text []byte
	text, err = os.ReadFile(inputPath)
	if err != nil {
		return encodeSize, encodeTime, fmt.Errorf("open input file %s failed: %v", inputPath, err.Error())
	}

	// get huffman codes
	var codes HuffmanCodes
	codes, err = GetHuffmanCodes(string(text))
	if err != nil {
		return encodeSize, encodeTime, fmt.Errorf("generate huffman codes failed: %v", err.Error())
	}
	var codeGenTime time.Time = time.Now()

	// create output directory and file
	var outputFile *os.File
	outputFile, err = OpenFile(outputPath)
	if err != nil {
		return encodeSize, encodeTime, fmt.Errorf("open output file %s failed: %v", outputPath, err.Error())
	}
	defer outputFile.Close()

	// write huffman table
	var huffmanTableSize int
	huffmanTableSize, err = writeHuffmanTable(outputFile, codes)
	if err != nil {
		return encodeSize, encodeTime, err
	}

	// write string
	var dataSize int
	dataSize, err = writeString(outputFile, text, codes)
	if err != nil {
		return encodeSize, encodeTime, err
	}
	var writeFileTime time.Time = time.Now()

	// write size and time record
	encodeSize = EncodeSize{
		orininal:     len(text),
		HuffmanTable: huffmanTableSize,
		EncodedData:  dataSize,
	}
	encodeTime = EncodeTime{
		CodeGenTime:   codeGenTime.Sub(startTime),
		WriteFileTime: writeFileTime.Sub(codeGenTime),
	}
	return encodeSize, encodeTime, err
}

func BatchEncode(inputPath string, outputPath string) (result BatchEncodeResult, err error) {
	// record start time
	var startTime time.Time = time.Now()
	var errors []BatchError = make([]BatchError, 0)

	// normalize input & output paths
	inputPath = filepath.Clean(inputPath)
	outputPath = filepath.Clean(outputPath)

	// collect input files
	var inputFiles []string
	var getFilesErrors []BatchError
	inputFiles, getFilesErrors, err = GetFilesInDir(inputPath)
	if err != nil {
		return result, fmt.Errorf("get input files failed: %v", err.Error())
	}
	errors = append(errors, getFilesErrors...)

	// get output paths for input files
	var outputPaths []string
	var getOutputPathErrors []BatchError
	outputPaths, getOutputPathErrors = GetOutputPaths(inputFiles, outputPath, "bin")
	errors = append(errors, getOutputPathErrors...)

	// process each file with goroutines
	var success int = 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, inputPath := range inputFiles {
		wg.Add(1)
		go func(idx int, inPath string) {
			defer wg.Done()
			// encode
			_, _, encErr := Encode(inPath, outputPaths[idx])
			mu.Lock()
			defer mu.Unlock()
			if encErr != nil {
				errors = append(errors, BatchError{Path: inPath, Err: encErr})
			} else {
				success++
			}
		}(i, inputPath)
	}

	wg.Wait()

	// fill result
	result = BatchEncodeResult{
		InputPath:    inputPath,
		OutputPath:   outputPath,
		TotalCount:   len(inputFiles),
		SuccessCount: success,
		Time:         time.Since(startTime),
		Errors:       errors,
	}
	return result, nil
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
func writeString(file io.Writer, text []byte, codes HuffmanCodes) (size int, err error) {
	// encode data
	var dataRecorder *BitsRecorder = NewBitsRecorder()
	// encoded data width (in bits)
	var dataWidth int64 = 0

	// write each byte to bits recorder
	for _, char := range text {
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
